package handlers

import (
	"encoding/json"
	"net/http"

	models "github.com/yhung-mea7/HeyNeighbor/account-service/data"
	"github.com/yhung-mea7/go-rest-kit/data"
)

func updateUser(repo models.IUserUpdate) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInformation := r.Context().Value(ak).(*models.User)
		requestBody := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {

			usrHandler.log.Println("[ERROR] unable to parse request body to map", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"Unable to process request body"}, rw)
			return
		}
		if err := repo.UpdateUser(userInformation.Username, requestBody); err != nil {
			usrHandler.log.Println("[ERROR] updating user: ", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"failed to update user"}, rw)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	}
}
