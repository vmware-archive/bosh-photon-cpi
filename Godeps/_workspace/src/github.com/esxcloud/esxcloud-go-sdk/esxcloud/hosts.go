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

// Creates a host.
func (api *HostsAPI) Create(hostSpec *HostCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(hostSpec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/hosts", bytes.NewBuffer(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Deletes a host with specified ID.
func (api *HostsAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/hosts/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Returns all hosts
func (api *HostsAPI) GetAll() (result *Hosts, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/hosts", api.client.options.Token)
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
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/hosts/"+id, api.client.options.Token)
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
	uri := api.client.Endpoint+"/v1/hosts/"+id+"/tasks"
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
