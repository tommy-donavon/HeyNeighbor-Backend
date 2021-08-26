package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/account-service/handlers"
)

func SetUpRoutes(sm *mux.Router, userHandler *handlers.UserHandler) {
	sm.Use(userHandler.GlobalContentTypeMiddleware)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/healthcheck", userHandler.HealthCheck())

	createAdminAccountRouter := sm.Methods(http.MethodPost).Subrouter()
	createAdminAccountRouter.Handle("/create-admin", userHandler.CreateAdminUser())
	createAdminAccountRouter.Use(userHandler.ValidateUserMiddleware)
}
