package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	jsonE "github.com/yhung-mea7/go-rest-kit/data"
)

func (ph *PropertyHandler) UpdateTenantInformation() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrCtx, err := ph.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			jsonE.ToJSON(&message{err}, rw)
			return
		}
		requestBody := map[string]interface{}{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			ph.log.Println("[ERROR] unable to parse request body to map", err)
			rw.WriteHeader(http.StatusBadRequest)
			jsonE.ToJSON(&message{"Unable to process request body"}, rw)
			return
		}
		usr := usrCtx.(*data.Tenant)
		if err := ph.repo.UpdateTenantInformation(usr.Username, requestBody); err != nil {
			ph.log.Println("[ERROR] error updating user", err)
			rw.WriteHeader(http.StatusBadRequest)
			jsonE.ToJSON(&message{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusAccepted)
	}
}
