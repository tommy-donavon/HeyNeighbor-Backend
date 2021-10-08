package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/maintenance-service/data"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

func updateMaintenanceRequest(repo data.IMaintenancePatch) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		requestBody := map[string]interface{}{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			myMainHandler.log.Println("[ERROR] unable to parse request body to map", err)
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{"Unable to process request body"}, rw)
			return
		}
		uCode, err := getRequestId(r)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		userCtx, err := myMainHandler.ctx.Get(r.Context(), "user")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		usr := userCtx.(*userInformation)
		if err := repo.UpdateMaintenanceRequest(uCode, usr.Username, requestBody); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}

		rw.WriteHeader(http.StatusAccepted)

	}
}
