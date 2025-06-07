package models

type User struct {
	UserID      string `gorm:"primaryKey"`
	UserName    string
	Faculty     string
	MachineID   string
	MachineName string
	OSVersion   string
}
