package weaviatedb

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

const ArticlesCollName = "articles"

const (
	articleTitleField       = "title"
	articleDescriptionField = "description"
	articleTopicsField      = "topics"
	articleCategoryField    = "category"
)

type ArticleSchema struct {
	weaviateClass *models.Class
	client        *weaviate.Client
}

func GetArticleSchema() *ArticleSchema {
	cls := &models.Class{
		Class:      ArticlesCollName,
		Vectorizer: "text2vec-openai", // Todo variabalize maybe
		ModuleConfig: map[string]interface{}{
			"text2vec-openai":   map[string]interface{}{},
			"generative-openai": map[string]interface{}{},
		},
		Properties: []*models.Property{
			{
				Name:     articleTitleField,
				DataType: []string{"text"},
			},
			{
				Name:     articleDescriptionField,
				DataType: []string{"text"},
			},
			{
				Name:     articleTopicsField,
				DataType: []string{"text[]"},
			},
			{
				Name:     articleCategoryField,
				DataType: []string{"text"},
			},
		},
	}
	return &ArticleSchema{
		weaviateClass: cls,
	}
}

func (a *ArticleSchema) SetWeaviateClient(c *weaviate.Client) {
	a.client = c
}

// TODO
func (a *ArticleSchema) Search(d any) (any, error) {
	return nil, nil
}

type NewArticle struct {
	ID          string
	Title       string
	Description string
	Topics      []string
	Category    string
}

func (a *ArticleSchema) Insert(d any) error {
	article := d.(NewArticle)
	object := &models.Object{
		Class: ArticlesCollName,
		Properties: map[string]any{
			articleTitleField:       article.Title,
			articleDescriptionField: article.Description,
			articleTopicsField:      article.Topics,
			articleCategoryField:    article.Category,
		},
	}

	if _, err := a.client.Batch().ObjectsBatcher().WithObjects(object).Do(context.Background()); err != nil {
		return err
	}
	return nil
}

func (a *ArticleSchema) GetSchema() *models.Class {
	return a.weaviateClass
}
