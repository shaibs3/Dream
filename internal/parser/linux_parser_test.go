package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinuxParser_Parse(t *testing.T) {
	sample := `USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.9  0.2 168332 11564 ?        Ss   01:10   0:01 /sbin/init splash
root         2  0.0  0.0      0     0 ?        S    01:10   0:00 [kthreadd]
root         3  0.0  0.0      0     0 ?        I<   01:10   0:00 [rcu_gp]
root         4  0.0  0.0      0     0 ?        I<   01:10   0:00 [rcu_par_gp]
root         5  0.0  0.0      0     0 ?        I<   01:10   0:00 [kworker/0:0H-kblockd]
root         6  0.0  0.0      0     0 ?        I<   01:10   0:00 [mm_percpu_wq]`

	parser := &LinuxParser{}
	entries, err := parser.Parse(sample)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(entries))

	first, ok := entries[0].(LinuxProcess)
	assert.True(t, ok, "First entry is not LinuxProcess")
	assert.Equal(t, "root", first.User)
	assert.Equal(t, 1, first.PID)
	assert.Equal(t, "/sbin/init splash", first.Command)
}
