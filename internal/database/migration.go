package database

import (
	"github.com/brettman/go-rest-api-course/internal/comment"
	"github.com/jinzhu/gorm"
)

// MigratedDB - creates comment table with a db migration
func MigrateDB(db *gorm.DB) error{
	if result := db.AutoMigrate(&comment.Comment{}); result.Error != nil{
		return result.Error
	}
	return nil
}