package categories

import "github.com/cblokkeel/newspaper/internal/db/mongo"

type Category struct {
    Name string `json:"name"`
    Topics []string `json:"topics"`
}

func CategoryFromModel(model *mongo.CategoryModel) Category {
    return Category{
        Name: model.Name,
        Topics: model.Topics,
    }
} 
