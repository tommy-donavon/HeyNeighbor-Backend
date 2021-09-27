package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/property-service/handlers"
)

func SetUpRoutes(sm *mux.Router, ph *handlers.PropertyHandler) {
	sm.Use(ph.GlobalContentTypeMiddleware)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/healthcheck", ph.HealthCheck())

	createPropertyRouter := sm.Methods(http.MethodPost).Subrouter()
	createPropertyRouter.Handle("/", ph.CreateProperty())
	createPropertyRouter.Use(ph.AuthMiddleware)
	createPropertyRouter.Use(ph.ValidatePropertyMiddleware)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.Handle("/tenant", ph.AddTenantToProperty())
	postRouter.Use(ph.AuthMiddleware)

	getProperties := sm.Methods(http.MethodGet).Subrouter()
	getProperties.Handle("/admin", ph.GetManagerProperties())
	getProperties.Handle("/code/{code}", ph.GetPropertyByServerCode())
	getProperties.Handle("/tenant", ph.GetTenantProperties())
	getProperties.Use(ph.AuthMiddleware)
}
