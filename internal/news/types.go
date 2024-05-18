package news

import (
	"time"

	"github.com/cblokkeel/newspaper/internal/db"
)

type Source struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Article struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Desc      string     `json:"desc"`
	Link      string     `json:"link"`
	Published *time.Time `json:"published"`
	Image     string     `json:"img"`
	Source    Source     `json:"source"`
	Upvotes   int        `json:"upvotes"`
	Downvotes int        `json:"downvotes"`
}

func SourceFromModel(model *db.SourceModel) Source {
	return Source{
		Name: model.Name,
		Icon: model.Icon,
	}
}

func ArticleFromModel(model *db.ArticleModel) Article {
	return Article{
		ID:        model.ID.Hex(),
		Title:     model.Title,
		Desc:      model.Desc,
		Link:      model.Link,
		Published: model.Published,
		Image:     model.Image,
		Source:    SourceFromModel(&model.Source),
		Upvotes:   model.Upvotes,
		Downvotes: model.Downvotes,
	}
}
