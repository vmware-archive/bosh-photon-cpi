package esxcloud

import (
	"bytes"
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
	"strconv"
)

type VmAPI struct {
	client *Client
}

type VmCreateSpec struct {
	Flavor        string         `json:"flavor"`
	SourceImageID string         `json:"sourceImageId"`
	AttachedDisks []AttachedDisk `json:"attachedDisks"`
	Affinities    []LocalitySpec `json:"affinities,omitempty"`
	Name          string         `json:"name"`
	Tags          []string       `json:"tags,omitempty"`
}

type VmOperation struct {
	Operation string                 `json:"operation"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type AttachedDisk struct {
	Flavor     string `json:"flavor"`
	Kind       string `json:"kind"`
	CapacityGB int    `json:"capacityGb"`
	Name       string `json:"name"`
	State      string `json:"state"`
	ID         string `json:"id,omitempty"`
	BootDisk   bool   `json:"bootDisk"`
}

type VM struct {
	SourceImageID string          `json:"sourceImageId,omitempty"`
	Cost          []QuotaLineItem `json:"cost"`
	Kind          string          `json:"kind"`
	AttachedDisks []AttachedDisk  `json:"attachedDisks"`
	Datastore     string          `json:"datastore,omitempty"`
	AttachedISOs  []ISO           `json:"attachedIsos,omitempty"`
	Tags          []string        `json:"tags,omitempty"`
	SelfLink      string          `json:"selfLink,omitempty"`
	Flavor        string          `json:"flavor"`
	Host          string          `json:"host,omitempty"`
	Name          string          `json:"name"`
	State         string          `json:"string"`
	ID            string          `json:"id"`
}

type ISO struct {
	Size int64  `json:"size,omitempty"`
	Kind string `json:"kind,omitempty"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

type VmDiskOperation struct {
	DiskID    string                 `json:"diskId"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

func (api *VmAPI) Get(id string) (vm *VM, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/vms/"+id)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	vm = &VM{}
	err = json.NewDecoder(res.Body).Decode(vm)
	return
}

func (api *VmAPI) Delete(id string, force bool) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/vms/"+id+"?force="+strconv.FormatBool(force))

	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) AttachDisk(id string, op *VmDiskOperation) (task *Task, err error) {
	body, err := json.Marshal(op)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/vms/"+id+"/attach_disk", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) DetachDisk(id string, op *VmDiskOperation) (task *Task, err error) {
	body, err := json.Marshal(op)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/vms/"+id+"/detach_disk", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) AttachISO(id, isoPath string) (task *Task, err error) {
	res, err := rest.MultipartUploadFile(api.client.httpClient, api.client.Endpoint+"/v1/vms/"+id+"/attach_iso", isoPath, nil)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := getTask(getError(res))
	return result, err
}

func (api *VmAPI) DetachISO(id string) (task *Task, err error) {
	body := []byte{}
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/vms/"+id+"/detach_iso", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) Operation(id string, op *VmOperation) (task *Task, err error) {
	body, err := json.Marshal(op)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/vms/"+id+"/operations", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
