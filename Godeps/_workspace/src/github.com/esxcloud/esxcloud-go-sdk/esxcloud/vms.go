package esxcloud

import (
	"bytes"
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
	"strconv"
)

// Contains functionality for VMs API.
type VmAPI struct {
	client *Client
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

func (api *VmAPI) SetMetadata(id string, metadata *VmMetadata) (task *Task, err error) {
	body, err := json.Marshal(metadata)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/vms/"+id+"/set_metadata", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
