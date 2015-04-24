package mocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"bytes"
	"io/ioutil"

	. "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
)

func NewMockServer() (server *httptest.Server) {
	return httptest.NewServer(
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, r.URL)
	}))
}

func NewMockTask(operation, state string, id string, steps ...Step) *Task {
	return &Task{Operation: operation, State: state, ID: id, Steps: steps}
}

func CreateResponder(response string) (Responder) {
	return Responder(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 0,
			Body: ioutil.NopCloser(bytes.NewBufferString(response)),
			ContentLength: int64(len(response)),
			Request: req,
		}

		resp.Header = make(map[string][]string)
		resp.Header.Add("Content-Type", "application/json")

		return resp, nil
	})
}