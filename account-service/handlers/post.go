package handlers

import (
	"net/http"

	models "github.com/yhung-mea7/HeyNeighbor/account-service/data"
	"github.com/yhung-mea7/go-rest-kit/data"
)

func loginUser(repo models.IUserRead) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrHandler.log.Println("POST LOGIN")
		login, ok := r.Context().Value(lk).(login)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to process login information"}, rw)
			return
		}
		user, err := repo.GetUser(login.Username)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&message{"Invalid Login information"}, rw)
			return
		}
		if err := models.CheckPassword(user.Password, login.Password); err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&message{"Invalid Login information"}, rw)
			return

		}
		token, err := usrHandler.jwt.CreateJwToken(user.Username, int(user.AccountType))
		if err != nil {
			usrHandler.log.Println(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"Failed to sign token"}, rw)
			return
		}
		data.ToJSON(&struct {
			Token string `json:"token"`
		}{token}, rw)

	}
}

func createUser(repo models.IUserCreate) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrHandler.log.Println("POST CREATE ADMIN USER")
		user, ok := r.Context().Value(uk).(models.User)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to process login information"}, rw)
			return
		}
		if err := repo.CreateUser(&user); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			usrHandler.log.Println(err)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
