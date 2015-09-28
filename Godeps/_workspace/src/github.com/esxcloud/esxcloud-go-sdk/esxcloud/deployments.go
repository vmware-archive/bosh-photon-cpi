package esxcloud

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

// Contains functionality for deployments API.
type DeploymentsAPI struct {
	client *Client
}

var DeploymentUrl string = "/deployments"

// Creates a deployment
func (api *DeploymentsAPI) Create(deploymentSpec *DeploymentCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(deploymentSpec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+DeploymentUrl, bytes.NewBuffer(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Deletes a deployment with specified ID.
func (api *DeploymentsAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+DeploymentUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Returns all deployments.
func (api *DeploymentsAPI) GetAll() (result *Deployments, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+DeploymentUrl, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &Deployments{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Gets a deployment with the specified ID.
func (api *DeploymentsAPI) Get(id string) (deployment *Deployment, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+DeploymentUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Deployment
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

// Gets all hosts with the specified deployment ID.
func (api *DeploymentsAPI) GetHosts(id string) (result *Hosts, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+DeploymentUrl+"/"+id+"/hosts", api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &Hosts{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Gets all the vms with the specified deployment ID.
func (api *DeploymentsAPI) GetVms(id string) (result *VMs, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+DeploymentUrl+"/"+id+"/vms", api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &VMs{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}
