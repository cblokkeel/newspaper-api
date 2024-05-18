package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDB(uri string, db string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(db)
	return &MongoDB{
		client,
		database,
	}, nil
}

func (m *MongoDB) getCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

func (m *MongoDB) Exists(ctx context.Context, collName string, filter bson.M) bool {
	coll := m.getCollection(collName)
	err := coll.FindOne(ctx, filter).Err()
	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func (m *MongoDB) Find(ctx context.Context, collName string, filter bson.M, opts *options.FindOptions) (*mongo.Cursor, error) {
    coll := m.getCollection(collName)
    return coll.Find(ctx, filter, opts)
}

func (m *MongoDB) Insert(ctx context.Context, collName string, obj interface{}) error {
	coll := m.getCollection(collName)
	_, err := coll.InsertOne(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Disconnect(ctx context.Context) {
	m.client.Disconnect(ctx)
}
