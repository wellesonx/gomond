package parser

import (
	"github.com/gelleson/gomond/gomond/pb"
	"strings"
	"time"
)

type Level int

func (l Level) ToString() string {
	switch l {
	case INFO:
		return "info"
	case DEBUG:
		return "debug"
	case WARNING:
		return "warning"
	case ERROR:
		return "error"
	case FATAL:
		return "fatal"
	default:
		return "unknown"
	}
}

func (l Level) ToPBLevel() pb.Level {
	switch l {
	case INFO:
		return pb.Level_info
	case DEBUG:
		return pb.Level_debug
	case WARNING:
		return pb.Level_warning
	case ERROR:
		return pb.Level_error
	case FATAL:
		return pb.Level_fatal
	default:
		return pb.Level_unknown
	}
}

type Type int

const (
	FILE Type = iota + 1
	HTTP
)

const (
	INFO Level = iota + 1
	DEBUG
	WARNING
	ERROR
	FATAL
	UNKNOWN
)

type Log struct {
	App       string    `json:"app"`
	Timestamp time.Time `json:"timestamp"`
	Level     Level     `json:"level"`
	Message   string    `json:"message"`
	File      string    `json:"file"`
	Type      Type      `json:"type"`
	Line      int64     `json:"line"`
	Payload   []byte    `json:"payload"`
}

func NewLog() *Log {
	return &Log{}
}

func stringToLevel(l string) Level {
	switch strings.ToLower(l) {
	case "info":
		return INFO
	case "debug":
		return DEBUG
	case "warning":
		return WARNING
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return UNKNOWN
	}
}
