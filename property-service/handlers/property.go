package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/property-service/data"

	"github.com/yhung-mea7/go-rest-kit/context"
	consul_register "github.com/yhung-mea7/go-rest-kit/register"
)

type (
	PropertyHandler struct {
		repo       *data.PropertyRepo
		log        *log.Logger
		validator  *data.Validation
		register   *consul_register.ConsulClient
		ctxHandler *context.ContextHandler
	}

	message struct {
		Message interface{} `json:"message"`
	}
)

//Creates new PropertyHandler.
func NewPropertyHandler(repo *data.PropertyRepo, log *log.Logger, v *data.Validation, register *consul_register.ConsulClient) *PropertyHandler {
	ph := &PropertyHandler{
		repo:      repo,
		log:       log,
		validator: v,
		register:  register,
	}
	ph.ctxHandler = &context.ContextHandler{}
	return ph
}

func getServerCode(r *http.Request) string {
	vars := mux.Vars(r)
	code, ok := vars["code"]
	if !ok {
		return ""
	}
	return code
}
