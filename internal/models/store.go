package models

import (
	"dream/internal/parser"

	"gorm.io/gorm"
)

func LinuxProcessFromParser(p parser.LinuxProcess) *LinuxProcess {
	return &LinuxProcess{
		User:    p.User,
		CPU:     p.CPU,
		MEM:     p.MEM,
		VSZ:     p.VSZ,
		RSS:     p.RSS,
		TTY:     p.TTY,
		Stat:    p.Stat,
		Start:   p.Start,
		Time:    p.Time,
		Command: p.Command,
	}
}

func WindowsProcessFromParser(p parser.WindowsProcess) *WindowsProcess {
	return &WindowsProcess{
		ImageName:   p.ImageName,
		SessionName: p.SessionName,
		SessionNum:  p.SessionNum,
		MemUsage:    p.MemUsage,
	}
}

func StoreMetadataAndProcesses(db *gorm.DB, meta *Metadata, processes []interface{}, osType string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(meta).Error; err != nil {
			return err
		}
		for _, p := range processes {
			switch osType {
			case "linux":
				parsed := p.(parser.LinuxProcess)
				linuxProc := LinuxProcessFromParser(parsed)
				process := &Process{
					MetadataID: meta.ID,
					PID:        parsed.PID,
					OSType:     "linux",
				}
				if err := tx.Create(process).Error; err != nil {
					return err
				}
				linuxProc.ProcessID = process.ID
				if err := tx.Create(linuxProc).Error; err != nil {
					return err
				}
			case "windows":
				parsed := p.(parser.WindowsProcess)
				windowsProc := WindowsProcessFromParser(parsed)
				process := &Process{
					MetadataID: meta.ID,
					PID:        parsed.PID,
					OSType:     "windows",
				}
				if err := tx.Create(process).Error; err != nil {
					return err
				}
				windowsProc.ProcessID = process.ID
				if err := tx.Create(windowsProc).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
