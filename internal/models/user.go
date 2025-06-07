package models

type User struct {
	ID          uint `gorm:"primaryKey"`
	UserID      string
	UserName    string
	Faculty     string
	MachineID   string
	MachineName string
	OSVersion   string
}
