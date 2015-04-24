package mocks

import (
	"encoding/json"
)

func ToJson(v interface{}) string {
	res, err := json.Marshal(v)
	if err != nil {
		// Since this method is only for testing, don't return
		// any errors, just panic.
		panic("Error serializing struct into JSON")
	}
	// json.Marshal returns []byte, convert to string
	return string(res[:])
}
