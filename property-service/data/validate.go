package data

import (
	"regexp"

	"github.com/go-playground/validator"
)

func (p *Property) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("zip", validateZipCode)
	return validate.Struct(p)
}

func validateZipCode(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`\d{5}(-\d{4})?`)
	return re.Match([]byte(fl.Field().String()))
}
