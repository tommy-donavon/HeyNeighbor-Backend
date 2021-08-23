package data

import (
	"regexp"

	"github.com/go-playground/validator"
)

func (u *User) Validate() error {
	validator := validator.New()
	validator.RegisterValidation("phone", validatePhoneNumber)
	return validator.Struct(u)
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^(\d{1,2}-)?(\d{3}-){2}\d{4}$`)
	return re.Match([]byte(fl.Field().String()))
}
