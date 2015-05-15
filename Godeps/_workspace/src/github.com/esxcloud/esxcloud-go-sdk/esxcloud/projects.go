package esxcloud

import (
	"bytes"
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

type ProjectsAPI struct {
	client *Client
}

type ProjectCreateSpec struct {
	ResourceTicket ResourceTicketReservation `json:"resourceTicket"`
	Name           string                    `json:"name"`
}

type ProjectList struct {
	Items []ProjectCompact `json:"items"`
}

type ProjectCompact struct {
	Kind           string        `json:"kind"`
	ResourceTicket ProjectTicket `json:"resourceTicket"`
	Name           string        `json:"name"`
	ID             string        `json:"id"`
	Tags           []string      `json:"tags"`
	SelfLink       string        `json:"selfLink"`
}

type ProjectTicket struct {
	TenantTicketID   string          `json:"tenantTicketId"`
	Usage            []QuotaLineItem `json:"usage"`
	TenantTicketName string          `json:"tenantTicketName"`
	Limits           []QuotaLineItem `json:"limits"`
}

func (api *ProjectsAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/projects/"+id)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *ProjectsAPI) CreateDisk(id string, spec *DiskCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/projects/"+id+"/disks", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *ProjectsAPI) FindDisks(id string, options *Options) (result *DiskList, err error) {
	uri := api.client.Endpoint + "/v1/projects/" + id + "/disks"
	if options != nil {
		uri += getQueryString(options)
	}
	res, err := rest.Get(api.client.httpClient, uri)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &DiskList{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

func (api *ProjectsAPI) CreateVM(id string, spec *VmCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/projects/"+id+"/vms", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
