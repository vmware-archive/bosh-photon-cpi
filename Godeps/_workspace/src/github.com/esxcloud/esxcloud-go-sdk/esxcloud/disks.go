package esxcloud

import (
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
	"strconv"
)

// Contains functionality for disks API.
type DisksAPI struct {
	client *Client
}

// Gets a PersistentDisk for the disk with specified ID.
func (api *DisksAPI) Get(diskID string) (disk *PersistentDisk, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/disks/"+diskID, api.client.options.Token)
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
func (api *DisksAPI) Delete(diskID string, force bool) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/disks/"+diskID+"?force="+strconv.FormatBool(force), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
