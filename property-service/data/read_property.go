package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type IPropertyRead interface {
	GetAllTenantProperties(string) ([]*Property, error)
	GetPropertyByServerCode(string) (*Property, error)
	GetAllManagerProperties(string) ([]*Property, error)
}

func (pr *PropertyRepo) GetAllTenantProperties(tenantUsername string) ([]*Property, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	result, err := coll.Find(
		ctx,
		bson.M{"tenants.username": tenantUsername},
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
