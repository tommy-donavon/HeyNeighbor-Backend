package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/property-service/register"
)

func main() {
	sm := mux.NewRouter()
	logger := log.New(os.Stdout, "property-service", log.LstdFlags)

	consulClient := register.NewConsulClient("property-service")
	consulClient.RegisterService()
	defer consulClient.DeregisterService()

	sm.Handle("/healthcheck", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		json.NewEncoder(rw).Encode(&struct {
			Message string `json:"message"`
		}{"gud"})
	}))
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
