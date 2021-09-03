package handlers

import (
	"log"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	"github.com/yhung-mea7/HeyNeighbor/property-service/register"
)

type (
	PropertyHandler struct {
		repo       *data.PropertyRepo
		log        *log.Logger
		validator  *data.Validation
		register   *register.ConsulClient
		ctxHandler *contextHandler
	}

	message struct {
		Message interface{} `json:"message"`
	}
	userInformation struct {
		Username   string `json:"username"`
		UserType   int    `json:"account_type"`
		ProfileUri string `json:"profile_uri"`
	}
)

//Creates new PropertyHandler.
func NewPropertyHandler(repo *data.PropertyRepo, log *log.Logger, v *data.Validation, register *register.ConsulClient) *PropertyHandler {
	ph := &PropertyHandler{
		repo:      repo,
		log:       log,
		validator: v,
		register:  register,
	}
	ph.ctxHandler = &contextHandler{}
	return ph
}
