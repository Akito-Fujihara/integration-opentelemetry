package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("database connect success!!")
	return db, nil
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln("database can't close")
	}
	log.Println("database close success!!")
}
