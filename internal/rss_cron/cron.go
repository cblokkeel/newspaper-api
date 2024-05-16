package rsscron

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/cblokkeel/newspaper/internal/db"
	"github.com/mmcdole/gofeed"
)

const feeds_coll = "feeds"

type RSSCron struct {
	mongo *db.MongoDB
}

func NewRSSCron(mongo *db.MongoDB) *RSSCron {
	return &RSSCron{
		mongo,
	}
}

func (c *RSSCron) Start(ctx context.Context) {
	cursor, err := c.mongo.Find(ctx, feeds_coll, bson.M{})
	if err != nil {
		log.Fatalf("Failed to fetch feeds: %+v", err)
	}
    defer cursor.Close(ctx)

    fp := gofeed.NewParser()

    for cursor.Next(ctx) {
        var feed db.Feed
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
            if (item.Image != nil) {
                imgLink = item.Image.URL
            }
            article := &db.Article{
                Title: item.Title,
                Desc: item.Description,
                Link: item.Link,
                Published: item.PublishedParsed,
                Image: imgLink, 
                Upvotes: 0,
                Downvotes: 0,
                Source: db.Source{
                    Name: feed.Name,
                    Icon: "todo",
                },
            }

            err := c.mongo.Insert(ctx, localizedColl, article)
            if err != nil {
                log.Printf("failed to insert new article %s: %+v\n", article.Title, err)
            }
        }
    }
}
