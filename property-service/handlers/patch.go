package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

func updateTenantInformation(repo data.IPropertyUpdate) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrCtx, err := instance.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{err}, rw)
			return
		}
		requestBody := map[string]interface{}{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			instance.log.Println("[ERROR] unable to parse request body to map", err)
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{"Unable to process request body"}, rw)
			return
		}
		usr := usrCtx.(*data.Tenant)
		if err := repo.UpdateTenantInformation(usr.Username, requestBody); err != nil {
			instance.log.Println("[ERROR] error updating user", err)
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusAccepted)
	}
}
