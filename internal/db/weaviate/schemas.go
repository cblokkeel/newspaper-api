package weaviatedb

import "github.com/weaviate/weaviate/entities/models"

type WeaviateSchema interface {
	GetSchema() *models.Class
	Insert(data any) error
	Search(data any) (any, error)
}
