package parser

type ParserType string

const (
	JSONParserType string = "json"
	TextParserType string = "text"
)

type Parser interface {
	Read([]byte) (Log, error)
}
