package handlers

import "net/http"

func (mh *MaintenanceHandler) Healthcheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

	}
}
