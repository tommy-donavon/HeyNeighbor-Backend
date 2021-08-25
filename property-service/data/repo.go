package data

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropertyRepo struct {
	client *mongo.Client
	dbName string
}

// Creates new mongo connection client
func NewPropertyRepo() *PropertyRepo {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		panic(err)
	}
	dbName := os.Getenv("MONGO_DB")
	coll := client.Database(dbName).Collection("properties")
	_, err = coll.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "address.street_address", Value: 1},
				{Key: "address.city", Value: 1},
				{Key: "address.state", Value: 1},
				{Key: "address.zip_code", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		panic(err)
	}
	return &PropertyRepo{client, dbName}
}
