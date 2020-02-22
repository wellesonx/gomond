package parser

type ParserType string

const (
	JSONParserType string = "json"
	TextParserType string = "text"
)

type LogParser interface {
	Read([]byte) (Log, error)
}
