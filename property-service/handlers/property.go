package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/yhung-mea7/go-rest-kit/context"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
	consul_register "github.com/yhung-mea7/go-rest-kit/register"
)

type (
	propertyHandler struct {
		log        *log.Logger
		validator  *my_json.Validation
		register   *consul_register.ConsulClient
		ctxHandler *context.ContextHandler
	}

	message struct {
		Message interface{} `json:"message"`
	}
)

var lock = &sync.Mutex{}
var instance *propertyHandler

//Creates new PropertyHandler.
func NewPropertyHandler(log *log.Logger, register *consul_register.ConsulClient) *propertyHandler {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			validator := my_json.NewValidator(
				my_json.ValidationOption{
					Name:      "zip",
					Operation: my_json.NewValidatorFunc(`^\d{5}(-\d{4})?$`),
				},
			)
			instance = &propertyHandler{
				log:        log,
				validator:  validator,
				register:   register,
				ctxHandler: &context.ContextHandler{},
			}

		}
	}
	return instance
}

func getServerCode(r *http.Request) string {
	vars := mux.Vars(r)
	code, ok := vars["code"]
	if !ok {
		return ""
	}
	return code
}
