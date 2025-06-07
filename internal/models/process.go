package models

import "time"

type LinuxProcess struct {
	ID                 uint      `gorm:"primaryKey"`
	UserID             string    // Foreign key to users table
	Timestamp          time.Time // When the process was seen
	User               string
	CPU                float64
	MEM                float64
	VSZ                int
	RSS                int
	TTY                string
	Stat               string
	Start              string
	Time               string
	Command            string
	FileSize           int64
	IsSystemExecutable bool
}

type WindowsProcess struct {
	ID                 uint      `gorm:"primaryKey"`
	UserID             string    // Foreign key to users table
	Timestamp          time.Time // When the process was seen
	ImageName          string
	SessionName        string
	SessionNum         int
	MemUsage           string
	FileSize           int64
	IsSystemExecutable bool
}
