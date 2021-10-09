package handlers

import (
	"net/http"

	models "github.com/yhung-mea7/HeyNeighbor/account-service/data"
	"github.com/yhung-mea7/go-rest-kit/data"
)

func healthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&message{"Service is good to go"}, rw)
	}
}

func getLoggedInUser(repo models.IUserRead) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		client, ok := r.Context().Value(ak).(*models.User)
		if !ok {
			usrHandler.log.Println("[ERROR] type assertion to data.User failed")
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"Unable to get user information"}, rw)
			return
		}
		usr, err := repo.GetUser(client.Username)
		if err != nil {
			usrHandler.log.Println("[ERROR] type assertion to data.User failed")
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"Unable to get user information"}, rw)
			return

		}
		usr.Password = ""
		data.ToJSON(&client, rw)
	}
}
