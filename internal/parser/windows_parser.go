package parser

import (
	"dream/internal/types"
	"strconv"
	"strings"
)

type WindowsProcess struct {
	ImageName   string
	PID         int
	SessionName string
	SessionNum  int
	MemUsage    string
}

type WindowsParser struct{}

func (p *WindowsParser) Parse(raw string) ([]types.ProcessEntry, error) {
	lines := strings.Split(raw, "\n")
	var entries []types.ProcessEntry
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "=") {
			continue
		}
		fields := strings.Fields(line)
		// Find the PID (first integer in the line)
		pidIdx := -1
		for i, f := range fields {
			if _, err := strconv.Atoi(f); err == nil {
				pidIdx = i
				break
			}
		}
		if pidIdx == -1 || len(fields) < pidIdx+4 {
			continue // skip malformed
		}
		imageName := strings.Join(fields[:pidIdx], " ")
		pid, _ := strconv.Atoi(fields[pidIdx])
		sessionName := fields[pidIdx+1]
		sessionNum, _ := strconv.Atoi(fields[pidIdx+2])
		memUsage := strings.Join(fields[pidIdx+3:], " ")
		entry := WindowsProcess{
			ImageName:   imageName,
			PID:         pid,
			SessionName: sessionName,
			SessionNum:  sessionNum,
			MemUsage:    memUsage,
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
