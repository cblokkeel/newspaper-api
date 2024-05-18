package news

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cblokkeel/newspaper/internal/db"
)

type NewsService struct {
	rdb   *redis.Client
	mongo *db.MongoDB
}

func NewNewsService(rdb *redis.Client, mongo *db.MongoDB) *NewsService {
	return &NewsService{
		rdb,
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
        var article db.ArticleModel
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
