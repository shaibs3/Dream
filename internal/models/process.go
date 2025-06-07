package models

type Process struct {
	ID         uint `gorm:"primaryKey"`
	MetadataID uint
	PID        int
	OSType     string
	Linux      *LinuxProcess   `gorm:"foreignKey:ProcessID"`
	Windows    *WindowsProcess `gorm:"foreignKey:ProcessID"`
}

type LinuxProcess struct {
	ProcessID uint `gorm:"primaryKey"`
	User      string
	CPU       float64
	MEM       float64
	VSZ       int
	RSS       int
	TTY       string
	Stat      string
	Start     string
	Time      string
	Command   string
}

type WindowsProcess struct {
	ProcessID   uint `gorm:"primaryKey"`
	ImageName   string
	SessionName string
	SessionNum  int
	MemUsage    string
}
