package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type IPropertyDelete interface {
	DeleteProperty(string) error
	RemoveChannelFromProperty(string, string) error
	RemoveTenantFromProperty(string, string) error
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

func (pr *PropertyRepo) RemoveChannelFromProperty(serverCode, channel string) error {
	if channel == general || channel == announcements || channel == events {
		return fmt.Errorf("%s is unable to be removed from list", channel)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	coll := pr.client.Database(pr.dbName).Collection("properties")
	defer cancel()
	prop, err := pr.GetPropertyByServerCode(serverCode)
	if err != nil {
		return err
	}
	inList := false
	for i, v := range prop.Channels {
		if v == channel {
			prop.Channels[i] = prop.Channels[len(prop.Tenants)-1]
			prop.Channels = prop.Channels[:len(prop.Channels)-1]
			inList = true
		}
	}
	if !inList {
		return fmt.Errorf("%s is not a channel in this property", channel)
	}
	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": prop.ID},
		bson.D{
			{
				Key:   "$set",
				Value: bson.D{{Key: "channels", Value: prop.Channels}},
			},
		},
	)
	return err
}
func (pr *PropertyRepo) RemoveTenantFromProperty(serverCode, tenUsername string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop, err := pr.GetPropertyByServerCode(serverCode)
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
