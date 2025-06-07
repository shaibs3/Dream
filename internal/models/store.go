package models

import (
	"dream/internal/parser"
	"time"

	"gorm.io/gorm"
)

// PostgresStorage implements ProcessStorage interface using GORM
type PostgresStorage struct {
	db *gorm.DB
}

// NewPostgresStorage creates a new PostgresStorage instance
func NewPostgresStorage(db *gorm.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

// StoreUserAndProcesses implements ProcessStorage interface
func (s *PostgresStorage) StoreUserAndProcesses(user interface{}, processes []interface{}, osType string, timestamp interface{}) error {
	u := user.(*User)
	t := timestamp.(time.Time)

	return s.db.Transaction(func(tx *gorm.DB) error {
		var existingUser User
		if err := tx.Where(User{UserID: u.UserID}).FirstOrCreate(&existingUser, *u).Error; err != nil {
			return err
		}
		for _, p := range processes {
			switch osType {
			case "linux":
				parsed := p.(parser.LinuxProcess)
				linuxProc := LinuxProcessFromParser(parsed, u.UserID, t)
				if err := tx.Create(linuxProc).Error; err != nil {
					return err
				}
			case "windows":
				parsed := p.(parser.WindowsProcess)
				windowsProc := WindowsProcessFromParser(parsed, u.UserID, t)
				if err := tx.Create(windowsProc).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func LinuxProcessFromParser(p parser.LinuxProcess, userID string, timestamp time.Time) *LinuxProcess {
	return &LinuxProcess{
		UserID:    userID,
		Timestamp: timestamp,
		User:      p.User,
		CPU:       p.CPU,
		MEM:       p.MEM,
		VSZ:       p.VSZ,
		RSS:       p.RSS,
		TTY:       p.TTY,
		Stat:      p.Stat,
		Start:     p.Start,
		Time:      p.Time,
		Command:   p.Command,
	}
}

func WindowsProcessFromParser(p parser.WindowsProcess, userID string, timestamp time.Time) *WindowsProcess {
	return &WindowsProcess{
		UserID:      userID,
		Timestamp:   timestamp,
		ImageName:   p.ImageName,
		SessionName: p.SessionName,
		SessionNum:  p.SessionNum,
		MemUsage:    p.MemUsage,
	}
}
