package parser

import (
	"dream/internal/types"
	"fmt"
	"log"
)

func GetParser(req types.MessageRequest) (types.Parser, error) {
	switch req.OSVersion {
	case "Windows 10":
		log.Printf("ParserFactory: Selected WindowsParser for OSVersion=%s", req.OSVersion)
		return &WindowsParser{}, nil
	case "Ubuntu", "Linux":
		log.Printf("ParserFactory: Selected LinuxParser for OSVersion=%s", req.OSVersion)
		return &LinuxParser{}, nil
	default:
		log.Printf("ParserFactory: Unsupported OS version: %s", req.OSVersion)
		return nil, fmt.Errorf("unsupported OS version: %s", req.OSVersion)
	}
}
