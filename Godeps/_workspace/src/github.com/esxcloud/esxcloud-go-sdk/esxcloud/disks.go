package esxcloud

import (
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
	"strconv"
)

type DisksAPI struct {
	client *Client
}

type LocalitySpec struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

type DiskCreateSpec struct {
	Flavor     string         `json:"flavor"`
	Kind       string         `json:"kind"`
	CapacityGB int            `json:"capacityGb"`
	Affinities []LocalitySpec `json:"localitySpec,omitempty"`
	Name       string         `json:"name"`
	Tags       []string       `json:"tags,omitempty"`
}

type PersistentDisk struct {
	Flavor     string          `json:"flavor"`
	Cost       []QuotaLineItem `json:"cost"`
	Kind       string          `json:"kind"`
	Datastore  string          `json:"datastore,omitempty"`
	CapacityGB int             `json:"capacityGb,omitempty"`
	Name       string          `json:"name"`
	State      string          `json:"state"`
	ID         string          `json:"id"`
	VMs        []string        `json:"vms"`
	Tags       []string        `json:"tags,omitempty"`
	SelfLink   string          `json:"selfLink,omitempty"`
}

type DiskList struct {
	Items []PersistentDisk `json:"items"`
}

type Options struct {
	Name string
}

func (api *DisksAPI) Get(id string) (disk *PersistentDisk, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/disks/"+id)
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

func (api *DisksAPI) Delete(id string, force bool) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/disks/"+id+"?force="+strconv.FormatBool(force))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
