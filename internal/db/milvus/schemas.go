package milvus

import "github.com/milvus-io/milvus-sdk-go/v2/entity"

const ArticlesCollName = "articles"

const (
	ArticleIDField          = "article_id"
	ArticlePublishDateField = "publish_date"
	ArticleVectorField      = "article_vector"
)

var articleSchema = &entity.Schema{
	CollectionName: ArticlesCollName,
	Description:    "News Article",
	Fields: []*entity.Field{
		{
			Name:       ArticleIDField,
			DataType:   entity.FieldTypeString,
			PrimaryKey: true,
			AutoID:     false,
		},
		{
			Name:       ArticlePublishDateField,
			DataType:   entity.FieldTypeFloat,
			PrimaryKey: false,
			AutoID:     false,
		},
		{
			Name:     ArticleVectorField,
			DataType: entity.FieldTypeFloatVector,
			TypeParams: map[string]string{
				"dim": "1024",
			},
		},
	},
}
