package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/account-service/data"
)

func (uh *UserHandler) LoginUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uh.log.Println("POST LOGIN")
		login, ok := r.Context().Value(lk).(login)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to process login information"}, rw)
			return
		}
		user, err := uh.repo.GetUser(login.Username)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&message{"Invalid Login information"}, rw)
			return
		}
		if err := data.CheckPassword(user.Password, login.Password); err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&message{"Invalid Login information"}, rw)
			return

		}
		token, err := uh.jwt.CreateJwToken(user.Username, int(user.AccountType))
		if err != nil {
			uh.log.Println(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"Failed to sign token"}, rw)
			return
		}
		data.ToJSON(&struct {
			Token string `json:"token"`
		}{token}, rw)

	}
}

func (uh *UserHandler) CreateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uh.log.Println("POST CREATE ADMIN USER")
		user, ok := r.Context().Value(uk).(data.User)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to process login information"}, rw)
			return
		}
		if err := uh.repo.CreateUser(&user); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			uh.log.Println(err)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
