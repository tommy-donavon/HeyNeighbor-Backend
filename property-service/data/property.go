package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	my_json "github.com/yhung-mea7/go-rest-kit/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	PropertyName    string             `bson:"property_name" json:"property_name" validate:"required"`
	PropertyAddress *Address           `bson:"address" json:"address" validate:"required,dive,required"`
	PropertyManager string             `bson:"property_manager" json:"property_manager"`
	NumOfUnits      uint               `bson:"num_of_units" json:"num_of_units" validate:"required"`
	ServerCode      string             `bson:"server_code" json:"server_code"`
	Tenants         []*Tenant          `bson:"tenants" json:"tenants" validate:"dive"`
	Channels        []string           `bson:"channels" json:"channels"`
}

//TODO add go routine to update user rent payments

//Inserts property document into mongo database
func (pr *PropertyRepo) CreateProperty(prop *Property) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop.Channels = append(prop.Channels, general)
	prop.Channels = append(prop.Channels, announcements)
	prop.Channels = append(prop.Channels, events)
	_, err := coll.InsertOne(ctx, &prop)
	return err
}

func (pr *PropertyRepo) GetPropertyByServerCode(serverCode string) (*Property, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	result := coll.FindOne(
		ctx,
		bson.M{"server_code": serverCode},
	)
	prop := Property{}
	if err := result.Decode(&prop); err != nil {
		return nil, err
	}
	return &prop, nil
}

func (pr *PropertyRepo) GetAllManagerProperties(propertyManager string) ([]*Property, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	result, err := coll.Find(
		ctx,
		bson.M{"property_manager": propertyManager},
	)
	if err != nil {
		return nil, err
	}
	results := []*Property{}
	for result.Next(ctx) {
		prop := Property{}
		if err := result.Decode(&prop); err != nil {
			return nil, err
		}
		results = append(results, &prop)
	}
	return results, nil
}

// updates base fields for a property
func (pr *PropertyRepo) UpdateProperty(server_code string, updateInfo map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	prop, err := pr.GetPropertyByServerCode(server_code)
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
	if err := my_json.NewValidator(my_json.ValidationOption{
		Name:      "zip",
		Operation: my_json.NewValidatorFunc(`^\d{5}(-\d{4})?$`),
	}).Validate(prop); err != nil {
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
func (pr *PropertyRepo) DeleteProperty(serverCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	prop, err := pr.GetPropertyByServerCode(serverCode)
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
