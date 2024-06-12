package mongo

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDB() (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(os.Getenv("MONGO_DB"))
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

func (m *MongoDB) Insert(ctx context.Context, collName string, obj interface{}) (interface{}, error) {
	coll := m.getCollection(collName)
	res, err := coll.InsertOne(ctx, obj)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (m *MongoDB) Disconnect(ctx context.Context) {
	m.client.Disconnect(ctx)
}
