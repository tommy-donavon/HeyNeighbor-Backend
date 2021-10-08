package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/go-rest-kit/context"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
	consul_register "github.com/yhung-mea7/go-rest-kit/register"
)

type (
	maintenanceHandler struct {
		log       *log.Logger
		validator *my_json.Validation
		register  *consul_register.ConsulClient
		ctx       *context.ContextHandler
	}
	message struct {
		Message interface{} `json:"message"`
	}
)

var lock = &sync.Mutex{}

var myMainHandler *maintenanceHandler

func NewMaintenanceHandler(log *log.Logger, reg *consul_register.ConsulClient) *maintenanceHandler {
	if myMainHandler == nil {
		lock.Lock()
		defer lock.Unlock()
		if myMainHandler == nil {

			validator := my_json.NewValidator()
			return &maintenanceHandler{
				log:       log,
				validator: validator,
				register:  reg,
				ctx:       &context.ContextHandler{},
			}
		}
	}
	return myMainHandler
}

func getServerCode(r *http.Request) string {
	vars := mux.Vars(r)
	code, ok := vars["code"]
	if !ok {
		return ""
	}
	return code
}

func getRequestId(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	code, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("please provide request id")
	}
	uCode, err := strconv.ParseUint(code, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(uCode), nil
}
