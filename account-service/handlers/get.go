package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/account-service/data"
)

func (uh *UserHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uh.log.Println("[GET] Healthcheck")
		data.ToJSON(&generalMessage{"Service is good to go"}, rw)
	}
}
