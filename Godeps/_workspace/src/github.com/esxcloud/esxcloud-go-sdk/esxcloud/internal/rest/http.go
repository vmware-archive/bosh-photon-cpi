package rest

import (
	"io"
	"net/http"
)

// Interface for abstracting away HTTP client implementation to enable testing
type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, bodyType string, body io.Reader) (resp *http.Response, err error)
	Delete(url string) (resp *http.Response, err error)
}

// Default, real implementation of HttpClient
type DefaultHttpClient struct{}

func (_ DefaultHttpClient) Get(uri string) (resp *http.Response, err error) {
	resp, err = http.Get(uri)
	return
}

func (_ DefaultHttpClient) Post(uri, bodyType string, body io.Reader) (resp *http.Response, err error) {
	resp, err = http.Post(uri, bodyType, body)
	return
}

func (_ DefaultHttpClient) Delete(uri string) (resp *http.Response, err error) {
	client := http.DefaultClient
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return
	}
	resp, err = client.Do(req)
	return
}
