package utils

import "fmt"

func GetRedisKey(category string, categoryID string, action string) string {
    return fmt.Sprintf("%s:%s:%s", category, categoryID, action)
}
