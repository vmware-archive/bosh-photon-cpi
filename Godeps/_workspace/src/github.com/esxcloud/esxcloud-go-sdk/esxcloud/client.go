package esxcloud

import (
	"net/http"
)

// Represents stateless context needed to call esxcloud APIs.
type Client struct {
	Endpoint   string
	httpClient *http.Client
	Status     *StatusAPI
	Tenants    *TenantsAPI
	Tasks      *TasksAPI
	Projects   *ProjectsAPI
	Flavors    *FlavorsAPI
	Images     *ImagesAPI
	Disks      *DisksAPI
	VMs        *VmAPI
}

func NewClient(endpoint string) (c *Client) {
	c = &Client{Endpoint: endpoint, httpClient: &http.Client{}}
	c.Status = &StatusAPI{c}
	c.Tenants = &TenantsAPI{c}
	c.Tasks = &TasksAPI{c}
	c.Projects = &ProjectsAPI{c}
	c.Flavors = &FlavorsAPI{c}
	c.Images = &ImagesAPI{c}
	c.Disks = &DisksAPI{c}
	c.VMs = &VmAPI{c}
	return
}

func NewTestClient(endpoint string, httpClient *http.Client) (c *Client) {
	c = &Client{Endpoint: endpoint, httpClient: httpClient}
	c.Status = &StatusAPI{c}
	c.Tenants = &TenantsAPI{c}
	c.Tasks = &TasksAPI{c}
	c.Projects = &ProjectsAPI{c}
	c.Flavors = &FlavorsAPI{c}
	c.Images = &ImagesAPI{c}
	c.Disks = &DisksAPI{c}
	c.VMs = &VmAPI{c}
	return
}
