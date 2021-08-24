package data

type Property struct {
	PropertyName          string    `json:"property_name" validate:"required"`
	PropertyStreetAddress string    `json:"property_street_address" validate:"required"`
	PropertyCity          string    `json:"property_city" validate:"required"`
	PropertyState         string    `json:"property_state" validate:"required"`
	PropertyZipCode       string    `json:"property_zip_code" validate:"required,zip"`
	PropertyManager       string    `json:"property_manager" validate:"required"`
	Tenets                []*string `json:"tenets"`
	Channels              []*string `json:"channels"`
}
