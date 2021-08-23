package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/tree/main/account-service/data"
)

type (
	UserHandler struct {
		repo *data.UserRepo
		log  *log.Logger
	}
	generalMessage struct {
		Message interface{} `json:"message"`
	}
	userInformation struct {
		Username string `json:"username"`
		UserType int    `json:"user_type"`
	}
	keyValue struct{}
)

func NewUserHandler(repo *data.UserRepo, log *log.Logger) *UserHandler {
	return &UserHandler{
		repo: repo,
		log:  log,
	}
}

func getUserId(r *http.Request) uint {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0
	}
	return uint(id)
}
