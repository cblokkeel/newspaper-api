package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedModel struct {
	ID       string   `bson:"_id"`
	Name     string   `bson:"name"`
	Desc     string   `bson:"desc"`
	Category string   `bson:"category"`
	Link     string   `bson:"link"`
	Topics   []string `bson:"topics"`
	Locale   string   `bson:"locale"`
}

type SourceModel struct {
	Name string `bson:"name"`
	Icon string `bson:"icon"`
}

type ArticleModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	Desc      string             `bson:"desc"`
	Link      string             `bson:"link"`
	Published *time.Time         `bson:"published"`
	Image     string             `bson:"img"`
	Source    SourceModel        `bson:"source"`
	Upvotes   int                `bson:"upvotes"`
	Downvotes int                `bson:"downvotes"`
	Category  string             `bson:"category"`
	Topics    []string           `bson:"topics"`
}

type CategoryModel struct {
	Name   string   `bson:"_id"`
	Topics []string `bson:"topics"`
	Locale string   `bson:"locale"`
}
