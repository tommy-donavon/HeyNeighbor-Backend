package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/maintenance-service/data"
)

func (mh *maintenanceHandler) SetupHandler(sm *mux.Router, repo *data.MaintenanceRepo) {
	sm.Handle("/healthcheck", healthcheck()).Methods(http.MethodGet)
	sm.Use(globalContentTypeMiddleware)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/{code:[A-z0-9-]+}", getPropertyMaintenanceRequest(repo))
	getRouter.Use(authMiddleware)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.Handle("/{code:[A-z0-9-]+}", createMaintenanceRequest(repo))
	postRouter.Use(authMiddleware)
	postRouter.Use(validateMaintenanceRequestMiddleware)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.Handle("/{id:[0-9]+}", deleteMaintenanceRequest(repo))
	deleteRouter.Use(authMiddleware)

	patchRouter := sm.Methods(http.MethodPatch).Subrouter()
	patchRouter.Handle("/{id:[0-9]+}", updateMaintenanceRequest(repo))
	patchRouter.Use(authMiddleware)
}
