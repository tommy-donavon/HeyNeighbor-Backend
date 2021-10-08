package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
)

func (ph *propertyHandler) SetUpRoutes(sm *mux.Router, repo *data.PropertyRepo) {
	sm.Use(globalContentTypeMiddleware)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/healthcheck", healthcheck())

	createPropertyRouter := sm.Methods(http.MethodPost).Subrouter()
	createPropertyRouter.Handle("/", createProperty(repo))
	createPropertyRouter.Use(authMiddleware)
	createPropertyRouter.Use(validatePropertyMiddleware)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.Handle("/tenant", addTenantToProperty(repo))
	postRouter.Use(authMiddleware)

	getProperties := sm.Methods(http.MethodGet).Subrouter()
	getProperties.Handle("/admin", getManagerProperties(repo))
	getProperties.Handle("/code/{code}", getPropertyByServerCode(repo))
	getProperties.Handle("/tenant", getTenantProperties(repo))
	getProperties.Use(authMiddleware)

	patchRouter := sm.Methods(http.MethodPatch).Subrouter()
	patchRouter.Handle("/tenant", updateTenantInformation(repo))
	patchRouter.Use(authMiddleware)
}
