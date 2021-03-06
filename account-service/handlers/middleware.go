package handlers

import (
	"context"
	"net/http"
	"strings"

	models "github.com/yhung-mea7/HeyNeighbor/account-service/data"
	"github.com/yhung-mea7/go-rest-kit/data"
)

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func globalContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func validateUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.User{}
		if err := data.FromJSON(&user, r.Body); err != nil {
			usrHandler.log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}
		if err := usrHandler.validator.Validate(user); err != nil {
			usrHandler.log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), uk, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func validateLoginInformation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		login := login{}
		if err := data.FromJSON(&login, r.Body); err != nil {
			usrHandler.log.Println("[ERROR] deserializing login", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}
		ctx := context.WithValue(r.Context(), lk, login)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			rw.WriteHeader(http.StatusForbidden)
			data.ToJSON(&message{"No token provided"}, rw)
			return
		}
		jwToken := strings.Split(token, "Bearer ")
		if len(jwToken) == 2 {
			token = strings.TrimSpace(jwToken[1])
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"Malformed Token"}, rw)
			return
		}
		claims, err := usrHandler.jwt.CheckToken(token)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&message{"Unauthorized"}, rw)
			return
		}
		client, err := usrHandler.repo.GetUser(claims.Username)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"can not find provied user"}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), ak, client)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}
