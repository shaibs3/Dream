package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Metadata{}, &Process{}, &LinuxProcess{}, &WindowsProcess{}); err != nil {
		return nil, err
	}
	return db, nil
}
