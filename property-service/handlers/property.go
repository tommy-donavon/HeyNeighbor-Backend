package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/HeyNeighbor/property-service/data"

	"github.com/yhung-mea7/go-rest-kit/context"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
	consul_register "github.com/yhung-mea7/go-rest-kit/register"
)

type (
	PropertyHandler struct {
		repo       *data.PropertyRepo
		log        *log.Logger
		validator  *my_json.Validation
		register   *consul_register.ConsulClient
		ctxHandler *context.ContextHandler
	}

	message struct {
		Message interface{} `json:"message"`
	}
)

//Creates new PropertyHandler.
func NewPropertyHandler(log *log.Logger, register *consul_register.ConsulClient) *PropertyHandler {
	repo := data.NewPropertyRepo()
	validator := my_json.NewValidator(
		my_json.ValidationOption{
			Name:      "zip",
			Operation: my_json.NewValidatorFunc(`^\d{5}(-\d{4})?$`),
		},
	)
	ph := &PropertyHandler{
		repo:      repo,
		log:       log,
		validator: validator,
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
