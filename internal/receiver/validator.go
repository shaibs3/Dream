package receiver

import (
	"dream/internal/types"
	"fmt"
)

type Validator interface {
	Validate(messageReq *types.MessageRequest) error
}

type MessageRequestValidator struct{}

func (v *MessageRequestValidator) Validate(messageReq *types.MessageRequest) error {
	if messageReq.MachineID == "" || messageReq.MachineName == "" || messageReq.CommandType == "" ||
		messageReq.UserName == "" || messageReq.UserID == "" {
		return fmt.Errorf("missing required fields")
	}
	return nil
}
