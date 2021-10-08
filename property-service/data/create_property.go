package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type IPropertyCreate interface {
	CreateProperty(*Property) error
	AddTenantToProperty(string, *Tenant) error
	AddChannelToProperty(servercode, channel string) error
}

//Inserts property document into mongo database
func (pr *PropertyRepo) CreateProperty(prop *Property) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop.ServerCode = generateServerCode(prop.PropertyName)
	prop.Channels = append(prop.Channels, general)
	prop.Channels = append(prop.Channels, announcements)
	prop.Channels = append(prop.Channels, events)
	_, err := coll.InsertOne(ctx, &prop)
	return err
}

// add tenant to property
func (pr *PropertyRepo) AddTenantToProperty(serverCode string, ten *Tenant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop, err := pr.GetPropertyByServerCode(serverCode)
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

func (pr *PropertyRepo) AddChannelToProperty(servercode, channel string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	coll := pr.client.Database(pr.dbName).Collection("properties")
	defer cancel()
	prop, err := pr.GetPropertyByServerCode(servercode)
	if err != nil {
		return err
	}
	for _, v := range prop.Channels {
		if v == channel {
			return fmt.Errorf("%s is already a channel in property", channel)
		}
	}
	prop.Channels = append(prop.Channels, channel)
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
