package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Tenant struct {
	Username   string `json:"username" validate:"required"`
	Nickname   string `json:"nickname"`
	UnitNumber uint   `json:"unit_number" validate:"required"`
}

//TODO add method to update tenant server nickname

// add tenant
func (pr *PropertyRepo) AddTenantToProperty(addr *Address, ten *Tenant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop, err := pr.GetProperty(addr)
	if err != nil {
		return err
	}
	for _, t := range prop.Tenants {
		if t.Username == ten.Username {
			return fmt.Errorf("%s is already in server", ten.Username)
		} else if t.Nickname == ten.Nickname {
			return fmt.Errorf("%s is already taken on server", ten.Nickname)
		}
	}
	prop.Tenants = append(prop.Tenants, ten)
	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": prop.ID},
		bson.D{
			{
				Key:   "$set",
				Value: bson.D{{Key: "tenants", Value: prop.Tenants}},
			},
		},
	)
	return err
}

func (pr *PropertyRepo) RemoveTenantFromProperty(addr *Address, tenUsername string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop, err := pr.GetProperty(addr)
	if err != nil {
		return err
	}
	inList := false
	for i, t := range prop.Tenants {
		if t.Username == tenUsername {
			prop.Tenants[i] = prop.Tenants[len(prop.Tenants)-1]
			prop.Tenants = prop.Tenants[:len(prop.Tenants)-1]
			inList = true
		}
	}
	if !inList {
		return fmt.Errorf("%s is not a tenant of this property", tenUsername)
	}
	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": prop.ID},
		bson.D{
			{
				Key:   "$set",
				Value: bson.D{{Key: "tenants", Value: prop.Tenants}},
			},
		},
	)
	return err

}
