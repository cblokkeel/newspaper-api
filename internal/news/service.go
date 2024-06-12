package news

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cblokkeel/newspaper/internal/db/mongo"
	redisdb "github.com/cblokkeel/newspaper/internal/db/redis"
)

type NewsService struct {
	redis   *redisdb.RedisDB
	mongo *mongo.MongoDB
}

func NewNewsService(redis *redisdb.RedisDB, mongo *mongo.MongoDB) *NewsService {
	return &NewsService{
		redis,
		mongo,
	}
}

func (s *NewsService) getNews(ctx context.Context) ([]Article, error) {
    var articles []Article
    // TODO handle pagination
    findOptions := options.Find()
    findOptions.SetLimit(20)
    cursor, err := s.mongo.Find(ctx, "fr_articles", bson.M{}, findOptions)
    if err != nil {
        return nil, err
    }

    for cursor.Next(ctx) {
        var article mongo.ArticleModel
        err := cursor.Decode(&article)
        if err != nil {
            // Decide what do
            return nil, err
        }
        articles = append(articles, ArticleFromModel(&article))
    }

    cursor.Close(ctx)
    return articles, nil
}
