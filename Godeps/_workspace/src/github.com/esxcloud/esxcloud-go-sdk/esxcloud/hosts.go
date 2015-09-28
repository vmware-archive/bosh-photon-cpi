package esxcloud

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

// Contains functionality for hosts API.
type HostsAPI struct {
	client *Client
}

var HostUrl string = "/hosts"

// Creates a host.
func (api *HostsAPI) Create(hostSpec *HostCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(hostSpec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+HostUrl, bytes.NewBuffer(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Deletes a host with specified ID.
func (api *HostsAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+HostUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Returns all hosts
func (api *HostsAPI) GetAll() (result *Hosts, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+HostUrl, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &Hosts{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Gets a host with the specified ID.
func (api *HostsAPI) Get(id string) (host *Host, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+HostUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Host
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

// Gets all tasks with the specified host ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *HostsAPI) GetTasks(id string, options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint + HostUrl + "/" + id + "/tasks"
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

// Gets all the vms with the specified deployment ID.
func (api *HostsAPI) GetVMs(id string) (result *VMs, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+HostUrl+"/"+id+"/vms", api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &VMs{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}
