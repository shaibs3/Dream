package receiver

import (
	"dream/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageRequestValidator_Validate(t *testing.T) {
	validator := &MessageRequestValidator{}

	// Valid case
	validReq := &types.MessageRequest{
		MachineID:   "123",
		MachineName: "TestMachine",
		OSVersion:   "Ubuntu",
		Timestamp:   types.MessageRequest{}.Timestamp,
		CommandType: "ps",
		UserName:    "student",
		UserID:      "u1",
		Faculty:     "Law",
	}
	assert.NoError(t, validator.Validate(validReq))

	// Missing MachineID
	invalidReq := *validReq
	invalidReq.MachineID = ""
	assert.Error(t, validator.Validate(&invalidReq))

	// Missing MachineName
	invalidReq = *validReq
	invalidReq.MachineName = ""
	assert.Error(t, validator.Validate(&invalidReq))

	// Missing CommandType
	invalidReq = *validReq
	invalidReq.CommandType = ""
	assert.Error(t, validator.Validate(&invalidReq))

	// Missing UserName
	invalidReq = *validReq
	invalidReq.UserName = ""
	assert.Error(t, validator.Validate(&invalidReq))

	// Missing UserID
	invalidReq = *validReq
	invalidReq.UserID = ""
	assert.Error(t, validator.Validate(&invalidReq))
}
