package esxcloud

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

// Contains functionality for tenants API.
type TenantsAPI struct {
	client *Client
}

// Options for GetResourceTickets API.
type ResourceTicketGetOptions struct {
	Name string `urlParam:"name"`
}

// Options for GetProjects API.
type ProjectGetOptions struct {
	Name string `urlParam:"name"`
}

var TenantUrl string = "/tenants"

// Returns all tenants on an esxcloud instance.
func (api *TenantsAPI) GetAll() (result *Tenants, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+TenantUrl, api.client.options.Token)
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

// Creates a tenant.
func (api *TenantsAPI) Create(tenantSpec *TenantCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(tenantSpec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+TenantUrl, bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Deletes the tenant with specified ID. Any projects, VMs, disks, etc., owned by the tenant must be deleted first.
func (api *TenantsAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+TenantUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Creates a resource ticket on the specified tenant.
func (api *TenantsAPI) CreateResourceTicket(tenantId string, spec *ResourceTicketCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+TenantUrl+"/"+tenantId+"/resource-tickets", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets resource tickets for tenant with the specified ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *TenantsAPI) GetResourceTickets(tenantId string, options *ResourceTicketGetOptions) (tickets *ResourceList, err error) {
	uri := api.client.Endpoint + TenantUrl + "/" + tenantId + "/resource-tickets"
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
	tickets = &ResourceList{}
	err = json.NewDecoder(res.Body).Decode(tickets)
	return
}

// Creates a project on the specified tenant.
func (api *TenantsAPI) CreateProject(tenantId string, spec *ProjectCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+TenantUrl+"/"+tenantId+"/projects", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets the projects for tenant with the specified ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *TenantsAPI) GetProjects(tenantId string, options *ProjectGetOptions) (result *ProjectList, err error) {
	uri := api.client.Endpoint + TenantUrl + "/" + tenantId + "/projects"
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
	result = &(ProjectList{})
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Gets all tasks with the specified tenant ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *TenantsAPI) GetTasks(id string, options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint + TenantUrl + "/" + id + "/tasks"
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
