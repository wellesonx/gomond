package parser

import (
	"github.com/araddon/dateparse"
	"github.com/juju/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	FILE_KEY    = "file"
	MESSAGE_KEY = "message"
	LINE_KEY    = "line"
	LEVEL_KEY   = "level"
	TIME_KEY    = "time"
)

const (
	FIRST_GROUP = 1
)

var (
	ERROR_LEVEL_NOT_FOUND     = errors.New("level extractor is not exist")
	ERROR_MESSAGE_NOT_FOUND   = errors.New("message extractor is not exist")
	ERROR_LEVEL_IS_REQUIRED   = errors.New("level extractor is not exist")
	ERROR_NOT_FOUND_SUBSTRING = errors.New("key substring is not found")
)

type TextOption struct {
	LevelField     string `json:"level"`
	MessageField   string `json:"message"`
	AppName        string `json:"app_name,omitempty"`
	LogName        string `json:"log_name,omitempty"`
	TimestampField string `json:"timestamp"`
	FileField      string `json:"file,omitempty"`
	LineField      string `json:"line,omitempty"`
}

type TextParser struct {
	config    TextOption
	kind      Type
	extractor map[string]regexp.Regexp
}

func NewTextParser(kind Type, config TextOption) (*TextParser, error) {

	extractor := make(map[string]regexp.Regexp)

	levelExtractor, err := regexp.Compile(config.LevelField)

	if err != nil {
		return &TextParser{}, errors.Annotate(err, "TextParser level regexp error")
	}

	if !isContainGroup(levelExtractor.SubexpNames(), "level") {
		return &TextParser{}, errors.New("there ins't Level name group in regex")
	}

	extractor[LEVEL_KEY] = *levelExtractor

	messageExtractor, err := regexp.Compile(config.MessageField)

	if err != nil {
		return &TextParser{}, errors.Annotate(err, "TextParser message regexp error")
	}

	if !isContainGroup(messageExtractor.SubexpNames(), "message") {
		return &TextParser{}, errors.New("there ins't Message name group in regex")
	}

	extractor[MESSAGE_KEY] = *messageExtractor

	fileExtractor, err := regexp.Compile(config.FileField)

	if err != nil {
		return &TextParser{}, errors.Annotate(err, "TextParser file regexp error")
	}

	if isContainGroup(fileExtractor.SubexpNames(), "file") {
		extractor[FILE_KEY] = *fileExtractor
	}

	timeExtractor, err := regexp.Compile(config.TimestampField)

	if err != nil {
		return &TextParser{}, errors.Annotate(err, "TextParser time regexp error")
	}

	if isContainGroup(timeExtractor.SubexpNames(), "time") {
		extractor[TIME_KEY] = *timeExtractor
	}

	lineExtractor, err := regexp.Compile(config.LineField)

	if err != nil {
		return &TextParser{}, errors.Annotate(err, "TextParser line regexp error")
	}

	if isContainGroup(lineExtractor.SubexpNames(), "line") {
		extractor[LINE_KEY] = *lineExtractor
	}

	return &TextParser{config: config, kind: kind, extractor: extractor}, nil
}

func (t TextParser) Parse(body []byte) (Log, error) {
	logObj := NewLog()
	record := string(body)

	levelEx, exist := t.extractor[LEVEL_KEY]

	if !exist {
		return Log{}, errors.Annotate(ERROR_LEVEL_NOT_FOUND, "TextParser parser")
	}

	levelSubmatch := levelEx.FindStringSubmatch(record)

	if len(levelSubmatch) <= 1 {
		return Log{}, errors.Annotate(ERROR_NOT_FOUND_SUBSTRING, "there isn't TextParser parser LEVEL substring in regex")
	}

	logObj.Level = stringToLevel(levelSubmatch[FIRST_GROUP])

	messageEx, exist := t.extractor[MESSAGE_KEY]

	if !exist {
		return Log{}, errors.Annotate(ERROR_MESSAGE_NOT_FOUND, "TextParser parser")
	}

	messageSubmatch := messageEx.FindStringSubmatch(record)

	if len(messageSubmatch) <= 1 {
		return Log{}, errors.Annotate(ERROR_NOT_FOUND_SUBSTRING, "there isn't TextParser parser MESSAGE substring in regex")
	}

	logObj.Message = messageSubmatch[FIRST_GROUP]

	timeEx, optionalExist := t.extractor[TIME_KEY]

	if !optionalExist {
		logObj.Timestamp = time.Now()
	} else {

		timeSubmatch := timeEx.FindStringSubmatch(record)

		if len(timeSubmatch) <= 1 {
			return Log{}, errors.Annotate(ERROR_NOT_FOUND_SUBSTRING, "there isn't TextParser parser TIME substring in regex")
		}

		parsedTime, err := dateparse.ParseAny(timeSubmatch[FIRST_GROUP])

		if err != nil {
			return Log{}, errors.Annotate(err, "TextParser parser parse datetime error")

		}
		logObj.Timestamp = parsedTime
	}

	fileEx, optionalExist := t.extractor[FILE_KEY]

	if optionalExist {
		fileSubmatch := fileEx.FindStringSubmatch(record)

		if len(fileSubmatch) <= 1 {
			return Log{}, errors.Annotate(ERROR_NOT_FOUND_SUBSTRING, "there isn't TextParser parser FILE substring in regex")
		}

		logObj.File = fileSubmatch[FIRST_GROUP]
	}

	lineEx, optionalExist := t.extractor[LINE_KEY]

	if optionalExist {
		lineSubmatch := lineEx.FindStringSubmatch(record)

		if len(lineSubmatch) <= 1 {
			return Log{}, errors.Annotate(ERROR_NOT_FOUND_SUBSTRING, "there isn't TextParser parser FILE substring in regex")
		}

		line, err := strconv.Atoi(lineSubmatch[FIRST_GROUP])

		if err != nil {
			return Log{}, errors.Annotate(err, "TextParser convert line from string to integer")
		}

		logObj.Line = int64(line)
	}

	logObj.Payload = body

	return *logObj, nil
}

func isContainGroup(groups []string, target string) bool {
	return strings.Contains(strings.Join(groups, " "), target)
}
