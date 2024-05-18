package deps

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/cblokkeel/newspaper/internal/categories"
	"github.com/cblokkeel/newspaper/internal/db"
	"github.com/cblokkeel/newspaper/internal/news"
)

var lock = &sync.Mutex{}

// DB
var (
	rdbInstance   *redis.Client
	mongoInstance *db.MongoDB
)

// NEWS
var (
	newsSvcInstance     *news.NewsService
	newsHandlerInstance *news.NewsHandler
)

// CATEGORIES

var (
	categoriesSvcInstance     *categories.CategoriesService
	categoriesHandlerInstance *categories.CategoriesHandler
)

func GetRDB() *redis.Client {
	if rdbInstance != nil {
		return rdbInstance
	}
	rdbInstance = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err := rdbInstance.Set(context.Background(), "foo", "bar", 0).Err(); err != nil {
		panic(err)
	}
	return rdbInstance
}

func GetMongo() *db.MongoDB {
	if mongoInstance != nil {
		return mongoInstance
	}
	m, err := db.NewMongoDB(os.Getenv("MONGO_URI"), os.Getenv("MONGO_DB"))
	if err != nil {
		panic(fmt.Errorf("Couldn't connect to mongo"))
	}
	mongoInstance = m
	return mongoInstance
}

func GetNewsSVC() *news.NewsService {
	if newsSvcInstance != nil {
		return newsSvcInstance
	}
	newsSvcInstance = news.NewNewsService(GetRDB(), GetMongo())
	return newsSvcInstance
}

func GetNewsHandler() *news.NewsHandler {
	if newsHandlerInstance != nil {
		return newsHandlerInstance
	}
	newsHandlerInstance = news.NewNewsHandler(GetNewsSVC())
	return newsHandlerInstance
}

func GetCategoriesSvc() *categories.CategoriesService {
    if categoriesSvcInstance != nil {
        return categoriesSvcInstance
    }
    categoriesSvcInstance := categories.NewCategoriesService(GetRDB(), GetMongo())
    return categoriesSvcInstance
}

func GetCategoriesHandler() *categories.CategoriesHandler {
	if categoriesHandlerInstance != nil {
		return categoriesHandlerInstance 
	}
    categoriesHandlerInstance = categories.NewCategoriesHandler(GetCategoriesSvc())
	return categoriesHandlerInstance
}
