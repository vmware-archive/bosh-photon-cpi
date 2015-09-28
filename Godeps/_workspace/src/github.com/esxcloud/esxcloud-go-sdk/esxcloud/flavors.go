package esxcloud

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

// Contains functionality for flavors API.
type FlavorsAPI struct {
	client *Client
}

// Options used for find/get APIs
type FlavorGetOptions struct {
	Name string `urlParam:"name"`
	Kind string `urlParam:"kind"`
}

var FlavorUrl string = "/flavors"

// Creates a flavor.
func (api *FlavorsAPI) Create(spec *FlavorCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+FlavorUrl, bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets details of flavor with specified ID.
func (api *FlavorsAPI) Get(flavorID string) (flavor *Flavor, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+FlavorUrl+"/"+flavorID, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	flavor = &Flavor{}
	err = json.NewDecoder(res.Body).Decode(flavor)
	return
}

// Gets flavors using options to filter results. Returns all flavors if options is nil.
func (api *FlavorsAPI) GetAll(options *FlavorGetOptions) (flavors *FlavorList, err error) {
	uri := api.client.Endpoint + FlavorUrl
	if options != nil {
		uri += getQueryString(options)
	}
	res, err := rest.Get(api.client.httpClient, uri, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	flavors = &FlavorList{}
	err = json.NewDecoder(res.Body).Decode(flavors)
	return
}

// Deletes flavor with specified ID.
func (api *FlavorsAPI) Delete(flavorID string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+FlavorUrl+"/"+flavorID, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets all tasks with the specified flavor ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *FlavorsAPI) GetTasks(id string, options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint+FlavorUrl+"/"+id+"/tasks"
	if options != nil {
		uri += getQueryString(options)
	}
	res, err := rest.Get(api.client.httpClient, uri, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &TaskList{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}
