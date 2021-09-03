package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

//base channels that every property needs to have
const (
	//general chat
	general string = "General"
	//landlord announcements
	announcements string = "Announcements"
	//events
	events string = "Events"
)

func (pr *PropertyRepo) AddChannelToProperty(addr *Address, channel string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	coll := pr.client.Database(pr.dbName).Collection("properties")
	defer cancel()
	prop, err := pr.GetProperty(addr)
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

func (pr *PropertyRepo) RemoveChannelFromProperty(addr *Address, channel string) error {
	if channel == general || channel == announcements || channel == events {
		return fmt.Errorf("%s is unable to be removed from list", channel)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	coll := pr.client.Database(pr.dbName).Collection("properties")
	defer cancel()
	prop, err := pr.GetProperty(addr)
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
