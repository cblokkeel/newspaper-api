package workers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/cblokkeel/newspaper/internal/db/mongo"
	redisdb "github.com/cblokkeel/newspaper/internal/db/redis"
	weaviatedb "github.com/cblokkeel/newspaper/internal/db/weaviate"
	"github.com/gofiber/fiber/v2"
	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RssWorker struct {
	CommonWorker

	mongo    *mongo.MongoDB
	weaviate *weaviatedb.WeaviateDB
}

// TODO externalize
const feeds_coll = "feeds"

func NewRssWorker(mongo *mongo.MongoDB, weaviate *weaviatedb.WeaviateDB, redis *redisdb.RedisDB, app *fiber.App) *RssWorker {
	w := &RssWorker{
		mongo:    mongo,
		weaviate: weaviate,
		CommonWorker: CommonWorker{
			name:  "rssworker",
			redis: redis,
		},
	}

	app.Post("/rss", func(c *fiber.Ctx) error {
		if locked := w.IsLocked(); locked {
			return c.Status(http.StatusConflict).SendString("RSS worker is locked")
		}
		go w.Start(context.Background())
		return c.SendString("RSS worker started")
	})

	return w
}

func (c *RssWorker) Periodicity() string {
	return "@every 1h"
}

func (c *RssWorker) Start(ctx context.Context) {
	if locked := c.IsLocked(); locked {
		log.Printf("worker %s is locked\n", c.name)
		return
	}
	if err := c.Lock(); err != nil {
		log.Printf("failed to lock rss worker")
		return
	}
	cursor, err := c.mongo.Find(ctx, feeds_coll, bson.M{}, options.Find())
	if err != nil {
		log.Fatalf("Failed to fetch feeds: %+v", err)
	}
	defer cursor.Close(ctx)

	fp := gofeed.NewParser()

	for cursor.Next(ctx) {
		var feed mongo.FeedModel
		err := cursor.Decode(&feed)
		if err != nil {
			log.Printf("failed to decode feed: %+v\n", err)
			continue
		}

		err = c.parseFeed(ctx, feed, fp)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}

	if err := c.Unlock(); err != nil {
		log.Printf("failed to unlock rss worker")
	}
}

func (c *RssWorker) parseFeed(ctx context.Context, feed mongo.FeedModel, parser *gofeed.Parser) error {
	log.Printf("Parsing feed %s\n", feed.Name)
	localizedColl := fmt.Sprintf("%s_articles", feed.Locale)
	feedParsed, err := parser.ParseURLWithContext(feed.Link, ctx) // TODO: add timeout
	if err != nil {
		return fmt.Errorf("failed to parse feed %s: %+v", feed.Link, err)
	}

	for _, item := range feedParsed.Items {
		if err := c.parseArticle(ctx, item, localizedColl, feed); err != nil {
			log.Println(err.Error())
		}
	}
	return nil
}

func (c *RssWorker) parseArticle(ctx context.Context, item *gofeed.Item, coll string, feed mongo.FeedModel) error {
	filter := bson.M{"link": item.Link}
	if exists := c.mongo.Exists(ctx, coll, filter); exists {
		return nil
	}

	articleObjectID, err := c.insertArticleInMongo(ctx, item, feed, coll)
	if err != nil {
		return err
	}

	if err := c.insertArticleInWeaviate(item, feed, articleObjectID); err != nil {
		return err
	}

	log.Printf("Successfully inserted article %s\n", item.Title)
	return nil
}

func (c *RssWorker) getImageLink(item *gofeed.Item) string {
	if item.Image != nil {
		return item.Image.URL
	}
	return ""
}

func (c *RssWorker) insertArticleInMongo(ctx context.Context, item *gofeed.Item, feed mongo.FeedModel, coll string) (*primitive.ObjectID, error) {
	article := &mongo.ArticleModel{
		ID:        primitive.NewObjectID(),
		Title:     item.Title,
		Desc:      item.Description,
		Link:      item.Link,
		Published: item.PublishedParsed,
		Image:     c.getImageLink(item),
		Upvotes:   0,
		Downvotes: 0,
		Source: mongo.SourceModel{
			Name: feed.Name,
			Icon: "todo",
		},
		Category: feed.Category,
		Topics:   feed.Topics,
	}

	articleID, err := c.mongo.Insert(ctx, coll, article)
	if err != nil {
		return nil, fmt.Errorf("failed to insert new article %s: %+v", article.Title, err)
	}

	articleObjectID, ok := articleID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to cast %v to object id", articleID)
	}

	return &articleObjectID, nil
}

func (c *RssWorker) insertArticleInWeaviate(item *gofeed.Item, feed mongo.FeedModel, mongoID *primitive.ObjectID) error {
	schema := c.weaviate.Schemas[weaviatedb.ArticlesCollName]
	if err := schema.Insert(weaviatedb.NewArticle{
		ID:          mongoID.Hex(),
		Title:       item.Title,
		Description: item.Description,
		Category:    feed.Category,
		Topics:      feed.Topics,
	}); err != nil {
		return fmt.Errorf("failed to insert article %s: %+v", item.Title, err)
	}
	return nil
}
