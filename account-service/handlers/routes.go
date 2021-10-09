package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/account-service/data"
)

func (uh *userHandler) SetUpRoutes(sm *mux.Router, repo *data.UserRepo) {
	sm.Use(globalContentTypeMiddleware)

	//get routers
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/healthcheck", healthCheck())

	authGet := sm.Methods(http.MethodGet).Subrouter()
	authGet.Handle("/", getLoggedInUser(repo))
	authGet.Use(authMiddleware)

	//post routers
	createAdminAccountRouter := sm.Methods(http.MethodPost).Subrouter()
	createAdminAccountRouter.Handle("/create-user", createUser(repo))
	createAdminAccountRouter.Use(validateUserMiddleware)

	logInAccountRouter := sm.Methods(http.MethodPost).Subrouter()
	logInAccountRouter.Handle("/", loginUser(repo))
	logInAccountRouter.Use(validateLoginInformation)

	//patch routers
	authPatch := sm.Methods(http.MethodPatch).Subrouter()
	authPatch.Handle("/", updateUser(repo))
	authPatch.Use(authMiddleware)

}
