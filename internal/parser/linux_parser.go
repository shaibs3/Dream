package parser

import (
	"dream/internal/types"
	"strconv"
	"strings"
)

type LinuxProcess struct {
	User    string
	PID     int
	CPU     float64
	MEM     float64
	VSZ     int
	RSS     int
	TTY     string
	Stat    string
	Start   string
	Time    string
	Command string
}

// LinuxParser parses ps auxww output
type LinuxParser struct{}

func (p *LinuxParser) Parse(raw string) ([]types.ProcessEntry, error) {
	lines := strings.Split(raw, "\n")
	if len(lines) < 2 {
		return nil, nil // or error
	}
	var entries []types.ProcessEntry
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue // skip malformed
		}
		pid, _ := strconv.Atoi(fields[1])
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		vsz, _ := strconv.Atoi(fields[4])
		rss, _ := strconv.Atoi(fields[5])
		entry := LinuxProcess{
			User:    fields[0],
			PID:     pid,
			CPU:     cpu,
			MEM:     mem,
			VSZ:     vsz,
			RSS:     rss,
			TTY:     fields[6],
			Stat:    fields[7],
			Start:   fields[8],
			Time:    fields[9],
			Command: strings.Join(fields[10:], " "),
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
