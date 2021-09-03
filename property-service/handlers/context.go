package handlers

import (
	"context"
	"errors"
	"fmt"
)

type (
	contextHandler struct{}

	contextValues struct {
		values map[string]interface{}
	}

	contextKey string
)

const (
	ck contextKey = "contextKey"
)

func (c *contextHandler) get(ctx context.Context, key string) (interface{}, error) {
	ctxMap, ok := ctx.Value(ck).(contextValues)
	if !ok {
		return nil, errors.New("unable to retrieve request context")
	}
	if value, ok := ctxMap.values[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("error retrieving item at key %s", key)
}

func (c *contextHandler) add(ctx context.Context, key string, value interface{}) context.Context {

	ctxMap, ok := ctx.Value(ck).(contextValues)
	if !ok {
		ctxMap = contextValues{
			values: make(map[string]interface{}),
		}
	}
	ctxMap.values[key] = value
	return context.WithValue(ctx, ck, ctxMap)

}
