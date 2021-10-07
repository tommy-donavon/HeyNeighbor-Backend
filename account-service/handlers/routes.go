package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetUpRoutes(sm *mux.Router, userHandler *UserHandler) {
	sm.Use(userHandler.GlobalContentTypeMiddleware)

	//get routers
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/healthcheck", userHandler.HealthCheck())

	authGet := sm.Methods(http.MethodGet).Subrouter()
	authGet.Handle("/", userHandler.GetLoggedInUser())
	authGet.Use(userHandler.Auth)

	//post routers
	createAdminAccountRouter := sm.Methods(http.MethodPost).Subrouter()
	createAdminAccountRouter.Handle("/create-user", userHandler.CreateUser())
	createAdminAccountRouter.Use(userHandler.ValidateUserMiddleware)

	logInAccountRouter := sm.Methods(http.MethodPost).Subrouter()
	logInAccountRouter.Handle("/", userHandler.LoginUser())
	logInAccountRouter.Use(userHandler.ValidateLoginInformation)

	//patch routers
	authPatch := sm.Methods(http.MethodPatch).Subrouter()
	authPatch.Handle("/", userHandler.UpdateUser())
	authPatch.Use(userHandler.Auth)

}
