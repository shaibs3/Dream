package parser

import (
	"dream/types"
	"fmt"
)

func GetParser(req types.MessageRequest) (types.Parser, error) {
	switch req.OSVersion {
	case "Windows 10":
		return &WindowsParser{}, nil
	case "Ubuntu", "Linux":
		return &LinuxParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported OS version: %s", req.OSVersion)
	}
}
