package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/maintenance-service/data"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

func deleteMaintenanceRequest(repo data.IMaintenanceDelete) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uCode, err := getRequestId(r)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		if err := repo.DeleteMaintenanceRequest(uint(uCode)); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusAccepted)
	}
}
