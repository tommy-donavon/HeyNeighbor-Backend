package handlers

import "log"

type (
	MaintenanceHandler struct {
		log *log.Logger
	}
)

func NewMaintenanceHandler(log *log.Logger) *MaintenanceHandler {
	return &MaintenanceHandler{
		log: log,
	}
}
