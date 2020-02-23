package parser

import (
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/gelleson/gomond/gomond/helper"
	"github.com/juju/errors"
	"time"
)

type JSONOption struct {
	LevelField     string `json:"level"`
	MessageField   string `json:"message"`
	AppName        string `json:"app_name,omitempty"`
	LogName        string `json:"log_name,omitempty"`
	TimestampField string `json:"timestamp"`
	FileField      string `json:"file,omitempty"`
}

type JSONParser struct {
	config JSONOption
	kind   Type
}

func NewJSONParser(kind Type, config JSONOption) (*JSONParser, error) {
	return &JSONParser{config: config, kind: kind}, nil
}

func (j JSONParser) Parse(body []byte) (Log, error) {
	logObject := NewLog()

	logObject.Type = j.kind

	logObject.App = j.config.AppName

	logObject.LogName = j.config.LogName

	fmt.Println(string(body))

	logObject.Payload = body

	parsed := make(helper.JSONBody)

	err := json.Unmarshal(body, &parsed)

	if err != nil {
		return Log{}, errors.Annotate(err, "JSONParser unmarshal")
	}

	if err = j.setLevel(parsed, logObject); err != nil {
		return Log{}, errors.Annotate(err, "JSONParser set level")
	}

	if err = j.setMessage(parsed, logObject); err != nil {
		return Log{}, errors.Annotate(err, "JSONParser set message")
	}

	if err = j.setTimestamp(parsed, logObject); err != nil {
		return Log{}, errors.Annotate(err, "JSONParser set timestamp")
	}

	j.setFile(parsed, logObject)

	return *logObject, nil
}

func (j JSONParser) setLevel(obj helper.JSONBody, logObj *Log) error {
	level, exist := obj.Get(j.config.LevelField)

	if !exist || level == "" {
		return errors.New("level field in JSONBody is not exist")
	}

	logObj.Level = stringToLevel(level)
	return nil
}

func (j JSONParser) setMessage(obj helper.JSONBody, logObj *Log) error {
	message, exist := obj.Get(j.config.MessageField)

	if !exist || message == "" {
		return errors.New("message field in JSONBody is not exist")
	}

	logObj.Message = message
	return nil
}

func (j JSONParser) setTimestamp(obj helper.JSONBody, logObj *Log) error {
	timestamp, exist := obj.Get(j.config.TimestampField)

	if !exist || timestamp == "" {
		logObj.Timestamp = time.Now()

		return nil
	}

	parsedTime, err := dateparse.ParseAny(timestamp)

	if err != nil {
		return err
	}

	logObj.Timestamp = parsedTime
	return nil
}

func (j JSONParser) setFile(obj helper.JSONBody, logObj *Log) {
	fileName, exist := obj.Get(j.config.FileField)

	if !exist || fileName == "" {
		return
	}

	logObj.File = fileName
}
