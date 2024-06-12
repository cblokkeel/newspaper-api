package weaviatedb

import "github.com/weaviate/weaviate/entities/models"

const ArticlesCollName = "articles"

const (
	articleTitleField       = "title"
	articleDescriptionField = "description"
	articleTopicsField      = "topics"
	articleCategoryField    = "category"
)

type ArticleSchema struct {
	weaviateClass *models.Class
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

// TODO
func (a *ArticleSchema) Search(d any) (any, error) {
	return nil, nil
}

// TODO
func (a *ArticleSchema) Insert(d any) error {
	return nil
}

func (a *ArticleSchema) GetSchema() *models.Class {
	return a.weaviateClass
}
