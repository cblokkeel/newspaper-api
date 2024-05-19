package db

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type MilvusDB struct {
	Client *client.Client
}

func NewMilvusDB(uri string) (*MilvusDB, error) {
	client, err := client.NewClient(context.Background(), client.Config{
		Address: uri,
	})
	if err != nil {
		return nil, err
	}
    exists, err := client.HasCollection(context.Background(), articlesCollName)
    if err != nil {
        return nil, err
    }

    if !exists {
        if client.CreateCollection(context.Background(), articleSchema, 2) != nil {
            return nil, err
        }
    }

	return &MilvusDB{
		Client: &client,
	}, nil
}

