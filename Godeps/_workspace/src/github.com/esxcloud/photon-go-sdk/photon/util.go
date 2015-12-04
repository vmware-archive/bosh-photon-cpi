package photon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Reads an error out of the HTTP response, or does nothing if
// no error occured.
func getError(res *http.Response) (*http.Response, error) {
	// Do nothing if the response is a successful 2xx
	if res.StatusCode/100 == 2 {
		return res, nil
	}
	var apiError ApiError
	// ReadAll is usually a bad practice, but here we need to read the response all
	// at once because we may attempt to use the data twice. It's preferable to use
	// methods that take io.Reader, e.g. json.NewDecoder
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &apiError)
	if err != nil {
		// If deserializing into ApiError fails, return a generic HttpError instead
		return nil, HttpError{res.StatusCode, string(body[:])}
	}
	apiError.HttpStatusCode = res.StatusCode
	return nil, apiError
}

// Reads a task object out of the HTTP response. Takes an error argument
// so that GetTask can easily wrap GetError. This function will do nothing
// if e is not nil.
// e.g. res, err := getTask(getError(someApi.Get()))
func getTask(res *http.Response, e error) (*Task, error) {
	if e != nil {
		return nil, e
	}
	var task Task
	err := json.NewDecoder(res.Body).Decode(&task)
	if err != nil {
		return nil, err
	}
	if task.State == "ERROR" {
		// Critical: return task as well, so that it can be examined
		// for error details.
		return &task, TaskError{task.ID}
	}
	return &task, nil
}

// Converts an options struct into a query string.
// E.g. type Foo struct {A int; B int} might return "?a=5&b=10".
// Will return an empty string if no options are set.
func getQueryString(options interface{}) string {
	buffer := bytes.Buffer{}
	buffer.WriteString("?")
	strct := reflect.ValueOf(options).Elem()
	typ := strct.Type()
	for i := 0; i < strct.NumField(); i++ {
		field := strct.Field(i)
		value := fmt.Sprint(field.Interface())
		if value != "" {
			buffer.WriteString(typ.Field(i).Tag.Get("urlParam") + "=" + url.QueryEscape(value))
			if i < strct.NumField()-1 {
				buffer.WriteString("&")
			}
		}
	}
	uri := buffer.String()
	if uri == "?" {
		return ""
	}
	return uri
}

// Sets security groups for a given entity (deployment/tenant/project)
func setSecurityGroups(client *Client, entityUrl string, securityGroups *SecurityGroups) (task *Task, err error) {
	body, err := json.Marshal(securityGroups)
	if err != nil {
		return
	}
	url := entityUrl + "/set_security_groups"
	res, err := rest.Post(client.httpClient, url, bytes.NewReader(body), client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
