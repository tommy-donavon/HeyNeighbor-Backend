package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/account-service/data"
)

func (uh *UserHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uh.log.Println("[GET] Healthcheck")
		data.ToJSON(&message{"Service is good to go"}, rw)
	}
}

func (uh *UserHandler) GetLoggedInUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		client, ok := r.Context().Value(ak).(*data.User)
		if !ok {
			uh.log.Println("[ERROR] type assertion to data.User failed")
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"Unable to get user information"}, rw)
			return
		}
		client.Password = ""
		data.ToJSON(&client, rw)
	}
}
