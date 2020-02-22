package parser

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TextParserSuite struct {
	suite.Suite
	option  TextOption
	message string
}

func (s *TextParserSuite) SetupTest() {
	s.option = TextOption{
		LevelField:     "\\[(?P<level>[\\w+]+)+\\]",
		MessageField:   "Desc:[ ]+(?P<message>[\\w+ ,\\d+.: ->]+)",
		AppName:        "app",
		TimestampField: "\\[[a-zA-Z]+\\](?P<time>[\\d+]{4}\\/[\\d+]{2}\\/[\\d+]{2} [\\d+]{2}:[\\d+]{2}:[\\d+]{2}.[\\d+]{3,7})",
		FileField:      "SrcFile:[ ]+(?P<file>[a-zA-Z.]+)",
		LineField:      "LineNum:[ ]+(?P<line>[0-9]+)",
	}

	message := `
	[Error]2020/02/21 14:33:18.055147
		SrcFile:  fxPxSession.go
		FuncName: github.com/quickfixgo/quickfix.(*fPSession).fxPxDownstreamMsgPostRoute
		LineNum:  1395
		Desc:     BackEnd connect successfully for FIXUser T42FC001,from 10.0.30.50:60850 -> 10.0.30.30:25142
	----------------------------------------------------------
	`
	s.message = message
}

func TestTextParser(t *testing.T) {

	s := new(TextParserSuite)

	suite.Run(t, s)
}