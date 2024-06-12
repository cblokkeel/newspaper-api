package weaviatedb

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type WeaviateDB struct {
	Schemas map[string]*WeaviateSchema
	c       *weaviate.Client
}

const (
	addr   = ""
	scheme = ""
)

func NewWeaviateDB(schemas []WeaviateSchema) *WeaviateDB {
	cfg := weaviate.Config{
		Host:   addr,
		Scheme: scheme,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	for _, schema := range schemas {
		client.Schema().ClassCreator().WithClass(schema.GetSchema()).Do(context.Background())
	}

	return &WeaviateDB{
		c: client,
	}
}
