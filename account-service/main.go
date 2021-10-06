package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/account-service/data"
	"github.com/yhung-mea7/HeyNeighbor/account-service/handlers"
	"github.com/yhung-mea7/HeyNeighbor/account-service/register"
	"github.com/yhung-mea7/HeyNeighbor/account-service/routes"
)

func main() {

	sm := mux.NewRouter()
	logger := log.New(os.Stdout, "account-service", log.LstdFlags)

	consulClient := register.NewConsulClient("account-service")
	consulClient.RegisterService() //register users-service to consul
	defer consulClient.DeregisterService()

	uh := handlers.NewUserHandler(data.NewUserRepo(), logger, os.Getenv("SECRET_KEY"))
	routes.SetUpRoutes(sm, uh)

	server := http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      sm,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Printf("Starting server on port: %v \n", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %v \n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	logger.Println("Got Signal:", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
