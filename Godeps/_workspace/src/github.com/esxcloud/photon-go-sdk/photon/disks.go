package photon

import (
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for disks API.
type DisksAPI struct {
	client *Client
}

var diskUrl string = "/disks/"

// Gets a PersistentDisk for the disk with specified ID.
func (api *DisksAPI) Get(diskID string) (disk *PersistentDisk, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+diskUrl+diskID, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	disk = &PersistentDisk{}
	err = json.NewDecoder(res.Body).Decode(disk)
	return
}

// Deletes a disk with the specified ID.
func (api *DisksAPI) Delete(diskID string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+diskUrl+diskID, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets all tasks with the specified disk ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *DisksAPI) GetTasks(id string, options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint + diskUrl + id + "/tasks"
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
