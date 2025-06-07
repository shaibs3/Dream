package parser

import (
	"testing"
)

func TestWindowsParser_Parse(t *testing.T) {
	sample := `Image Name                     PID Session Name        Session#    Mem Usage
========================= ======== ================ =========== ============
System Idle Process              0 Services                   0         24 K
System                           4 Services                   0     43,064 K
smss.exe                       400 Services                   0      1,548 K
csrss.exe                      564 Services                   0      6,144 K
wininit.exe                    652 Services                   0      5,044 K
csrss.exe                      676 Console                    1      9,392 K`

	parser := &WindowsParser{}
	entries, err := parser.Parse(sample)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(entries) != 6 {
		t.Errorf("Expected 6 entries, got %d", len(entries))
	}

	first, ok := entries[0].(WindowsProcess)
	if !ok {
		t.Fatalf("First entry is not WindowsProcess")
	}
	if first.ImageName != "System" || first.PID != 4 || first.MemUsage != "43,064 K" {
		t.Errorf("Unexpected first entry: %+v", first)
	}
}
