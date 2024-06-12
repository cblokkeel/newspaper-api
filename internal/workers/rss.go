package workers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cblokkeel/newspaper/internal/db/mongo"
	weaviatedb "github.com/cblokkeel/newspaper/internal/db/weaviate"
	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RssWorker struct {
	mongo    *mongo.MongoDB
	weaviate *weaviatedb.WeaviateDB
}

// TODO externalize
const feeds_coll = "feeds"

func NewRssWorker(mongo *mongo.MongoDB, weaviate *weaviatedb.WeaviateDB) *RssWorker {
	return &RssWorker{
		mongo,
		weaviate,
	}
}

func (c *RssWorker) Periodicity() string {
	return "@every 1h"
}

func (c *RssWorker) Start(ctx context.Context) {
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
		log.Printf("Parsing feed %s\n", feed.Name)
		localizedColl := fmt.Sprintf("%s_articles", feed.Locale)
		feedParsed, err := fp.ParseURLWithContext(feed.Link, ctx) // TODO: add timeout
		if err != nil {
			log.Printf("failed to parse feed %s: %+v\n", feed.Link, err)
			continue
		}

		for _, item := range feedParsed.Items {
			filter := bson.M{"link": item.Link}
			if exists := c.mongo.Exists(ctx, localizedColl, filter); exists {
				continue
			}
			var imgLink string
			if item.Image != nil {
				imgLink = item.Image.URL
			}
			article := &mongo.ArticleModel{
				ID:        primitive.NewObjectID(),
				Title:     item.Title,
				Desc:      item.Description,
				Link:      item.Link,
				Published: item.PublishedParsed,
				Image:     imgLink,
				Upvotes:   0,
				Downvotes: 0,
				Source: mongo.SourceModel{
					Name: feed.Name,
					Icon: "todo",
				},
				Category: feed.Category,
				Topics:   feed.Topics,
			}

			articleID, err := c.mongo.Insert(ctx, localizedColl, article)
			if err != nil {
				log.Printf("failed to insert new article %s: %+v\n", article.Title, err)
				continue
			}

			articleObjectID, ok := articleID.(primitive.ObjectID)
			if !ok {
				log.Printf("failed to cast %v to object id\n", articleObjectID)
				continue
			}

			schema := c.weaviate.Schemas[weaviatedb.ArticlesCollName]
			if err := schema.Insert(fmt.Sprintf("article title: %s; description: %s, category: %s, topics: %s", article.Title, article.Desc, article.Category, strings.Join(article.Topics, ","))); err != nil {
				log.Printf("Failed to insert article %s: %+v\n", article.Title, err)
				continue
			}
			log.Printf("Successfully inserted article %s\n", article.Title)
		}
	}
}
