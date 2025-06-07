package models

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	maxAttempts := 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("InitDB: failed to connect to database (attempt %d/%d): %v", attempt, maxAttempts, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&User{}, &LinuxProcess{}, &WindowsProcess{}); err != nil {
		return nil, err
	}
	return db, nil
}
