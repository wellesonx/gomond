package parser

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type JSONParserSuite struct {
	suite.Suite
	jp      *JSONParser
	message []byte
}

func (s *JSONParserSuite) SetupTest() {
	message := `{
				  "level": "info",
				  "message": "its ok",
				  "file": "domain.go",
				  "line": "1443",
				  "timestamp": "2020-02-12T02:14:15"
				}`

	s.message = []byte(message)

	option := JSONOption{
		LevelField:     "level",
		MessageField:   "message",
		AppName:        "watchers",
		TimestampField: "timestamp",
		FileField:      "file",
	}

	parser, _ := NewJSONParser(FILE, option)
	s.jp = parser
}

func (s JSONParserSuite) TestRead() {
	logObj, err := s.jp.Parse(s.message)

	s.Assert().Nil(err)

	s.Assert().Equal(s.jp.config.AppName, logObj.App)

	message := `{
				  "level": "info",
				  "message": "its ok",
				  "line": "1443",
				  "timestamp": "2020-02-12T02:14:15"
				}`

	logObj, err = s.jp.Parse([]byte(message))

	s.Assert().Nil(err)

	s.Assert().Empty(logObj.File)

	message = `{
				  "level": "info",
				  "message": "its ok",
				  "line": "1443",
				  "timestamp": "2020-02-12T02:14:15"
				}`

	logObj, err = s.jp.Parse([]byte(message))

	s.Assert().Nil(err)

	s.Assert().Empty(logObj.File)

	message = `{
				  "level": "info",
				  "line": "1443",
				  "timestamp": "2020-02-12T02:14:15"
				}`

	logObj, err = s.jp.Parse([]byte(message))

	s.Assert().Error(err)

	s.Assert().Empty(logObj)

	message = `{
				  "line": "1443",
				  "timestamp": "2020-02-12T02:14:15"
				}`

	logObj, err = s.jp.Parse([]byte(message))

	s.Assert().Error(err)

	s.Assert().Empty(logObj)

	message = `{
				  "level": "info",
				  "message": "its ok",
				  "line": "1443"
				}`

	logObj, err = s.jp.Parse([]byte(message))

	s.Assert().Nil(err, err)

	message = `{
				  "level": "info",
				  "message": "its ok",
				  "line": "1443",
				  "timestamp": "20200212T021415"
				}`

	logObj, err = s.jp.Parse([]byte(message))

	s.Assert().Error(err, err)

	message = `{
				  "level": "info",
				  "message": "its ok",
				  "line": "1443",
				  "timestamp": "20200212T021415",
				}`

	logObj, err = s.jp.Parse([]byte(message))

	s.Assert().Error(err, err)

}

func TestJSONPArser(t *testing.T) {
	s := new(JSONParserSuite)

	suite.Run(t, s)
}
