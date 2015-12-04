package photon

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for VMs API.
type VmAPI struct {
	client *Client
}

var vmUrl string = "/vms/"

func (api *VmAPI) Get(id string) (vm *VM, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+vmUrl+id, api.client.options.Token)
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

func (api *VmAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+vmUrl+id, api.client.options.Token)

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
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/attach_disk", bytes.NewReader(body), api.client.options.Token)
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
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/detach_disk", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) AttachISO(id, isoPath string) (task *Task, err error) {
	res, err := rest.MultipartUploadFile(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/attach_iso", isoPath, nil, api.client.options.Token)
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
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/detach_iso", bytes.NewReader(body), api.client.options.Token)
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
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/operations", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) Start(id string) (task *Task, err error) {
	body := []byte{}
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/start", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) Stop(id string) (task *Task, err error) {
	body := []byte{}
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/stop", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) Restart(id string) (task *Task, err error) {
	body := []byte{}
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/restart", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) Resume(id string) (task *Task, err error) {
	body := []byte{}
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/resume", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) Suspend(id string) (task *Task, err error) {
	body := []byte{}
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/suspend", bytes.NewReader(body), api.client.options.Token)
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
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/set_metadata", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets all tasks with the specified vm ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *VmAPI) GetTasks(id string, options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint + vmUrl + id + "/tasks"
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

func (api *VmAPI) GetNetworks(id string) (task *Task, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/networks", api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) GetMKSTicket(id string) (task *Task, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/mks_ticket", api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *VmAPI) SetTag(id string, tag *VmTag) (task *Task, err error) {
	body, err := json.Marshal(tag)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+vmUrl+id+"/tags", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
