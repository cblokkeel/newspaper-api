package weaviatedb

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

type WeaviateSchema interface {
	GetSchema() *models.Class
	Insert(data any) error
	Search(data any) (any, error)
	SetWeaviateClient(*weaviate.Client)
}
