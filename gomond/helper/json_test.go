package helper

import "testing"

func TestJSONBody_Get(t *testing.T) {
	body := make(JSONBody)

	body["case"] = "case1"

	case1, exist := body.Get("case")

	if !exist {
		t.Fatal("case is not exist")
	}

	if case1 != "case1" {
		t.Fatal("case is not exist")
	}

	case2, exist := body.Get("case2")

	if exist {
		t.Fatal("case2 is  exist")
	}

	if case2 != "" {
		t.Fatal("case is not exist")
	}
}
