package handlers

import (
	"log"
	"os"
	"sync"

	"github.com/yhung-mea7/HeyNeighbor/account-service/auth"
	"github.com/yhung-mea7/HeyNeighbor/account-service/data"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

type (
	userHandler struct {
		repo      *data.UserRepo
		log       *log.Logger
		jwt       *auth.JwtWrapper
		validator *my_json.Validation
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

var (
	lock       = &sync.Mutex{}
	usrHandler *userHandler
)

func NewUserHandler(log *log.Logger) *userHandler {
	if usrHandler == nil {
		lock.Lock()
		defer lock.Unlock()
		if usrHandler == nil {
			repo := data.NewUserRepo()
			validator := my_json.NewValidator(
				my_json.ValidationOption{
					Name:      "phone",
					Operation: my_json.NewValidatorFunc(`^(\d{1,2}-)?(\d{3}-){2}\d{4}$`),
				},
			)
			usrHandler = &userHandler{
				repo: repo,
				log:  log,
				jwt: &auth.JwtWrapper{
					SecretKey:       os.Getenv("SECRET_KEY"),
					Issuer:          "account-service",
					ExpirationHours: 24,
				},
				validator: validator,
			}
		}
	}
	return usrHandler
}
