package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(sm *mux.Router, mainRouter *MaintenanceHandler) {
	sm.Handle("/healthcheck", mainRouter.Healthcheck()).Methods(http.MethodGet)
}
