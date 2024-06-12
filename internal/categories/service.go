package categories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cblokkeel/newspaper/internal/db/mongo"
	redisdb "github.com/cblokkeel/newspaper/internal/db/redis"
)

type CategoriesService struct {
	redis   *redisdb.RedisDB
	mongo *mongo.MongoDB
}

func NewCategoriesService(redis *redisdb.RedisDB, mongo *mongo.MongoDB) *CategoriesService {
	return &CategoriesService{
		redis,
		mongo,
	}
}

func (s *CategoriesService) getCategories(ctx context.Context, locale string) ([]Category, error) {
	var categories []Category
	cursor, err := s.mongo.Find(ctx, "categories", bson.M{
		"locale": locale,
	}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}
    for cursor.Next(ctx) {
        var cat mongo.CategoryModel
        err := cursor.Decode(&cat)
        if err != nil {
            // Decide what do
            return nil, err
        }
        categories = append(categories, CategoryFromModel(&cat))
    }
    cursor.Close(ctx)
    return categories, nil
}
