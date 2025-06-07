package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, 6, len(entries))

	first, ok := entries[0].(WindowsProcess)
	assert.True(t, ok, "First entry is not WindowsProcess")
	assert.Equal(t, "System Idle Process", first.ImageName)
	assert.Equal(t, 0, first.PID)
	assert.Equal(t, "24 K", first.MemUsage)

	second, ok := entries[1].(WindowsProcess)
	assert.True(t, ok, "Second entry is not WindowsProcess")
	assert.Equal(t, "System", second.ImageName)
	assert.Equal(t, 4, second.PID)
	assert.Equal(t, "43,064 K", second.MemUsage)
}
