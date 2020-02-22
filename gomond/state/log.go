package state

import "github.com/gelleson/gomond/gomond/pb"

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
)

type Log struct {
	App     string `json:"app"`
	Level   Level  `json:"level"`
	Message string `json:"message"`
	File    string `json:"file"`
	Type    Type   `json:"type"`
	Line    int64  `json:"line"`
	Payload []byte `json:"payload"`
}

func NewLog() *Log {
	return &Log{}
}
