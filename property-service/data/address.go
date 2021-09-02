package data

type Address struct {
	StreetAddress string `json:"street_address" validate:"required"`
	City          string `json:"city" validate:"required"`
	State         string `json:"state" validate:"required"`
	ZipCode       string `json:"zip_code" validate:"required,zip"`
}
