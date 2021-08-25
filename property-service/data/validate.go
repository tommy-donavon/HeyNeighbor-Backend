package data

import (
	"regexp"

	"github.com/go-playground/validator"
)

type Validation struct {
	validate *validator.Validate
}

func NewValidator() *Validation {
	validate := validator.New()
	validate.RegisterValidation("zip", validateZipCode)
	return &Validation{validate}
}

func (v *Validation) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

func validateZipCode(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^\d{5}(-\d{4})?$`)
	return re.Match([]byte(fl.Field().String()))
}
