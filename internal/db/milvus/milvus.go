package milvus

import (
	"context"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type MilvusDB struct {
	Client client.Client
}

func NewMilvusDB(uri string) (MilvusDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	milvus, err := client.NewClient(ctx, client.Config{
		Address: uri,
	})
	if err != nil {
		return MilvusDB{}, err
	}
	exists, err := milvus.HasCollection(context.Background(), ArticlesCollName)
	if err != nil {
		return MilvusDB{}, err
	}

	if !exists {
		if milvus.CreateCollection(context.Background(), articleSchema, 2) != nil {
		    return MilvusDB{}, err
		}
	}

    return MilvusDB{
		Client: milvus,
    }, nil
}
