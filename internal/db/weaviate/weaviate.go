package weaviatedb

import (
	"context"
	"log"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type WeaviateDB struct {
	Schemas map[string]WeaviateSchema
	c       *weaviate.Client
}

const (
	addr   = "localhost:8081"
	scheme = "http"
)

func NewWeaviateDB(schemas []WeaviateSchema) *WeaviateDB {
	cfg := weaviate.Config{
		Host:   addr,
		Scheme: scheme,
	}

	client, _ := weaviate.NewClient(cfg)
	if err := client.WaitForWeavaite(time.Second * 10); err != nil {
		log.Fatalf("Weaviate is not available: %v", err)
	}

	for _, schema := range schemas {
		client.Schema().ClassCreator().WithClass(schema.GetSchema()).Do(context.Background())
		schema.SetWeaviateClient(client)
	}

	return &WeaviateDB{
		c: client,
	}
}
