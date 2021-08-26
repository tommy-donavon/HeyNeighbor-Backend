package handlers

import (
	"log"

	"github.com/yhung-mea7/HeyNeighbor/account-service/auth"
	"github.com/yhung-mea7/HeyNeighbor/account-service/data"
)

type (
	UserHandler struct {
		repo *data.UserRepo
		log  *log.Logger
		jwt  *auth.JwtWrapper
	}
	message struct {
		Message interface{} `json:"message"`
	}

	contextKey string
)

const (
	uk contextKey = "userKey"
	lk contextKey = "loginKey"
	ak contextKey = "adminKey"
)

func NewUserHandler(repo *data.UserRepo, log *log.Logger, key string) *UserHandler {
	return &UserHandler{
		repo: repo,
		log:  log,
		jwt: &auth.JwtWrapper{
			SecretKey:       key,
			Issuer:          "account-service",
			ExpirationHours: 24,
		},
	}
}

// func getUserName(r *http.Request) uint {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["username"])
// 	if err != nil {
// 		return 0
// 	}
// 	return uint(id)
// }
