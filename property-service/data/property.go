package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	PropertyName    string             `json:"property_name" validate:"required"`
	PropertyAddress *Address           `json:"address" validate:"required,dive,required"`
	PropertyManager string             `json:"property_manager"`
	NumOfUnits      uint               `json:"num_of_units" validate:"required"`
	ServerCode      string             `json:"server_code"`
	Tenants         []*Tenant          `json:"tenants" validate:"dive"`
	Channels        []string           `json:"channels"`
}

//Inserts property document into mongo database
func (pr *PropertyRepo) CreateProperty(prop *Property) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop.Channels = append(prop.Channels, general)
	prop.Channels = append(prop.Channels, announcements)
	prop.Channels = append(prop.Channels, events)
	_, err := coll.InsertOne(ctx, prop)
	return err
}

// queries property collection by provided address.
//
// returns a reference to found property or an error
func (pr *PropertyRepo) GetProperty(addr *Address) (*Property, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	result := coll.FindOne(ctx, bson.M{
		"address.street_address": addr.StreetAddress,
		"address.city":           addr.City,
		"address.state":          addr.State,
		"address.zip_code":       addr.ZipCode,
	})
	prop := Property{}
	if err := result.Decode(prop); err != nil {
		return nil, err
	}
	return &prop, nil
}

// updates base fields for a property
func (pr *PropertyRepo) UpdateProperty(addr *Address, updateInfo map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	prop, err := pr.GetProperty(addr)
	if err != nil {
		return err
	}
	propBytes, err := json.Marshal(prop)
	if err != nil {
		return err
	}
	propMap := map[string]interface{}{}
	err = json.Unmarshal(propBytes, &propMap)
	if err != nil {
		return err
	}

	for key, value := range updateInfo {
		if _, ok := propMap[key]; ok {
			switch key {
			case "property_name":
				v, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				prop.PropertyName = v
			case "address":
				v, ok := value.(Address)
				if !ok {
					return fmt.Errorf("error asserting value to address")
				}
				err = NewValidator().Validate(v)
				if err != nil {
					return err
				}
				prop.PropertyAddress = &v
			case "property_manager":
				v, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				prop.PropertyManager = v
			case "num_of_units":
				v, ok := value.(uint)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				prop.NumOfUnits = v
			}
		}
	}
	err = NewValidator().Validate(prop)
	if err != nil {
		return err
	}

	coll := pr.client.Database(pr.dbName).Collection("properties")
	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": prop.ID},
		bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "property_name", Value: prop.PropertyName},
					{Key: "address", Value: prop.PropertyAddress},
					{Key: "property_manager", Value: prop.PropertyManager},
					{Key: "num_of_units", Value: prop.NumOfUnits},
				},
			},
		},
	)
	return err
}

// deletes a property
func (pr *PropertyRepo) DeleteProperty(addr *Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	prop, err := pr.GetProperty(addr)
	if err != nil {
		return err
	}
	coll := pr.client.Database(pr.dbName).Collection("properties")
	_, err = coll.DeleteOne(
		ctx,
		bson.M{"_id": prop.ID},
	)
	return err
}
