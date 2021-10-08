package data

import (
	"context"
	"math/rand"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropertyRepo struct {
	client *mongo.Client
	dbName string
}

var lock = &sync.Mutex{}
var instance *PropertyRepo

// Creates new mongo connection client
func NewPropertyRepo() *PropertyRepo {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			rand.Seed(time.Now().UnixNano())
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
			indexes := []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "address.street_address", Value: 1},
						{Key: "address.city", Value: 1},
						{Key: "address.state", Value: 1},
						{Key: "address.zip_code", Value: 1},
					},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys:    bson.D{{Key: "server_code", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
			}
			_, err = coll.Indexes().CreateMany(ctx, indexes)
			if err != nil {
				panic(err)
			}
			instance = &PropertyRepo{client, dbName}
		}
	}
	return instance
}

func generateServerCode(propName string) string {
	b := make([]rune, 10)
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return propName + "-" + string(b)
}
