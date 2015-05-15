package esxcloud

import (
	"bytes"
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
	"net/url"
)

type TenantsAPI struct {
	client *Client
}

type Tenant struct {
	Projects        []BaseCompact `json:"projects"`
	ResourceTickets []BaseCompact `json:"resourceTickets"`
	Kind            string        `json:"kind"`
	Name            string        `json:"name"`
	ID              string        `json:"id"`
	SelfLink        string        `json:"selfLink"`
	Tags            []string      `json:"tags"`
}

type Tenants struct {
	Items []Tenant `json:"items"`
}

type TenantCreateSpec struct {
	Name string `json:"name"`
}

func (api *TenantsAPI) GetAll() (result *Tenants, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/tenants")
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &Tenants{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

func (api *TenantsAPI) Create(tenantSpec *TenantCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(tenantSpec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/tenants", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *TenantsAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/tenants/"+id)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *TenantsAPI) CreateResourceTicket(tenantId string, spec *ResourceTicketCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/tenants/"+tenantId+"/resource-tickets", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *TenantsAPI) GetResourceTickets(tenantId string, name *string) (tickets *ResourceList, err error) {
	uri := api.client.Endpoint + "/v1/tenants/" + tenantId + "/resource-tickets"
	if name != nil {
		uri += "?name=" + url.QueryEscape(*name)
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
	tickets = &ResourceList{}
	err = json.NewDecoder(res.Body).Decode(tickets)
	return
}

func (api *TenantsAPI) CreateProject(tenantId string, spec *ProjectCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/tenants/"+tenantId+"/projects", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *TenantsAPI) GetProjects(tenantId string, name *string) (result *ProjectList, err error) {
	uri := api.client.Endpoint + "/v1/tenants/" + tenantId + "/projects"
	if name != nil {
		uri += "?name=" + url.QueryEscape(*name)
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
	result = &(ProjectList{})
	err = json.NewDecoder(res.Body).Decode(result)
	return
}
