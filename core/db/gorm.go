package db

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLConnection() (*gorm.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) {
	// db.AutoMigrate(&entity.User{})
}
