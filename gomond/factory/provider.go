package factory

import (
	"encoding/json"
	"fmt"
	"github.com/gelleson/gomond/gomond/provider"
	"github.com/juju/errors"
)

func Provider(kind string, config *json.RawMessage) (provider.Provider, error) {

	switch kind {
	case "file":
		option := provider.FileOption{}

		err := json.Unmarshal(*config, &option)

		if err != nil {
			return nil, errors.Annotate(err, "Provider factory option read error")
		}

		fileProvider := provider.NewFileProvider(option)

		return fileProvider, nil

	default:
		return nil, errors.New(fmt.Sprintf("%s is not supported as provider", kind))
	}
}
