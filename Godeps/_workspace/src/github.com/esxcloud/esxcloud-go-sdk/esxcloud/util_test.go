package esxcloud

import (
	"testing"
)

type options struct {
	A int
	B string
}

func TestGetQueryString(t *testing.T) {
	opts := &options{5, "a test"}
	query := getQueryString(opts)
	if query != "?a=5&b=a+test" {
		t.Error("Query string is not correct")
	}
}
