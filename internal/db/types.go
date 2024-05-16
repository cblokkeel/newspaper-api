package db

import "time"

type Feed struct {
	ID       string   `bson:"_id"`
	Name     string   `bson:"name"`
	Desc     string   `bson:"desc"`
	Category string   `bson:"category"`
	Link     string   `bson:"link"`
	Topics   []string `bson:"topics"`
	Locale   string   `bson:"locale"`
}

type Source struct {
	Name string `bson:"name"`
	Icon string `bson:"icon"`
}

type Article struct {
	Title     string     `bson:"title"`
	Desc      string     `bson:"desc"`
	Link      string     `bson:"link"`
	Published *time.Time `bson:"published"`
	Image     string     `bson:"img"`
	Source    Source     `bson:"source"`
	Upvotes   int        `bson:"upvotes"`
	Downvotes int        `bson:"downvotes"`
}
