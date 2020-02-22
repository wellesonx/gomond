package helper

import "fmt"

type JSONBody map[string]interface{}

func (b JSONBody) Get(key string) (string, bool) {
	value, exist := b[key]

	if !exist {
		return "", false
	}

	return fmt.Sprintf("%v", value), true
}
