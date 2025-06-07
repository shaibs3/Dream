package models

import "time"

type Metadata struct {
	ID          uint `gorm:"primaryKey"`
	MachineID   string
	MachineName string
	OSVersion   string
	Timestamp   time.Time
	CommandType string
	BlobURL     string
	UserName    string
	UserID      string
	Processes   []Process `gorm:"foreignKey:MetadataID"`
}
