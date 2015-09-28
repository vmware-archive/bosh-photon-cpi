package esxcloud

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

// Contains functionality for networks API.
type NetworksAPI struct {
	client *Client
}

// Options used for GetAll API
type NetworkGetOptions struct {
	Name string `urlParam:"name"`
}

var NetworkUrl string = "/networks"

// Creates a network.
func (api *NetworksAPI) Create(networkSpec *NetworkCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(networkSpec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+NetworkUrl, bytes.NewBuffer(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Deletes a network with specified ID.
func (api *NetworksAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+NetworkUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets a network with the specified ID.
func (api *NetworksAPI) Get(id string) (network *Network, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+NetworkUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Network
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

// Returns all networks
func (api *NetworksAPI) GetAll(options *NetworkGetOptions) (result *Networks, err error) {
	uri := api.client.Endpoint + NetworkUrl
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
	result = &Networks{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}
