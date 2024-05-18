package categories

import "github.com/cblokkeel/newspaper/internal/db"

type Category struct {
    Name string `json:"name"`
    Topics []string `json:"topics"`
    Locale string `json:"locale"`
}

func CategoryFromModel(model *db.CategoryModel) Category {
    return Category{
        Name: model.Name,
        Topics: model.Topics,
        Locale: model.Locale,
    }
} 
