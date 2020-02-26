package factory

import (
	"encoding/json"
	"github.com/gelleson/gomond/gomond/parser"
	"github.com/juju/errors"
)

func Parser(app, name string, kind string, config *json.RawMessage) (parser.Parser, error) {
	switch kind {
	case parser.JSONParserType:
		jsonOption := parser.JSONOption{}

		jsonOption.AppName = app
		jsonOption.LogName = name

		err := json.Unmarshal(*config, &jsonOption)

		if err != nil {
			return nil, errors.Annotate(err, "Parse factory JSON option reading error")
		}

		p, _ := parser.NewJSONParser(parser.FILE, jsonOption)
		return p, nil

	case parser.TextParserType:
		textOption := parser.TextOption{}

		err := json.Unmarshal(*config, &textOption)

		if err != nil {
			return nil, errors.Annotate(err, "Parse factory TEXT option reading error")
		}

		textOption.AppName = app
		textOption.LogName = name

		p, err := parser.NewTextParser(parser.FILE, textOption)

		if err != nil {
			return nil, errors.Annotate(err, "Parse factory TEXT option validation error")
		}

		return p, nil

	default:
		return nil, errors.New("Unsupported type of parser")
	}

}
