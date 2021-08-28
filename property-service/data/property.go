package data

import (
	"context"
	"time"
)

type (
	Property struct {
		PropertyName    string    `json:"property_name" validate:"required"`
		PropertyAddress *Address  `json:"address" validate:"required,dive,required"`
		PropertyManager string    `json:"property_manager" validate:"required"`
		NumOfUnits      uint      `json:"num_of_units" validate:"required"`
		Tenants         []*string `json:"tenants"`
		Channels        []*string `json:"channels"`
	}
	Address struct {
		StreetAddress string `json:"street_address" validate:"required"`
		City          string `json:"city" validate:"required"`
		State         string `json:"state" validate:"required"`
		ZipCode       string `json:"zip_code" validate:"required,zip"`
	}
)

//Inserts property document into mongo database
func (pr *PropertyRepo) CreateProperty(prop *Property) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	_, err := coll.InsertOne(ctx, prop)
	return err
}
