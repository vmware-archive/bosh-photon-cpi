package photon

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for deployments API.
type DeploymentsAPI struct {
	client *Client
}

var deploymentUrl string = "/deployments"

// Creates a deployment
func (api *DeploymentsAPI) Create(deploymentSpec *DeploymentCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(deploymentSpec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+deploymentUrl, bytes.NewBuffer(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Deletes a deployment with specified ID.
func (api *DeploymentsAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+deploymentUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Deploys a deployment with specified ID.
func (api *DeploymentsAPI) Deploy(id string) (task *Task, err error) {
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+deploymentUrl+"/"+id+"/deploy", bytes.NewBuffer([]byte("")), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Destroys a deployment with specified ID.
func (api *DeploymentsAPI) Destroy(id string) (task *Task, err error) {
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+deploymentUrl+"/"+id+"/destroy", bytes.NewBuffer([]byte("")), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Returns all deployments.
func (api *DeploymentsAPI) GetAll() (result *Deployments, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+deploymentUrl, api.client.options.Token)
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
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+deploymentUrl+"/"+id, api.client.options.Token)
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
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+deploymentUrl+"/"+id+"/hosts", api.client.options.Token)
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
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+deploymentUrl+"/"+id+"/vms", api.client.options.Token)
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

// Initialize deployment migration from source to destination
func (api *DeploymentsAPI) InitializeDeploymentMigration(sourceDeploymentId string, id string) (task *Task, err error) {
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+deploymentUrl+"/"+id+"/initialize_migration", bytes.NewBuffer([]byte(id)), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Finalize deployment migration from source to destination
func (api *DeploymentsAPI) FinalizeDeploymentMigration(sourceDeploymentId string, id string) (task *Task, err error) {
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+deploymentUrl+"/"+id+"/finalize_migration", bytes.NewBuffer([]byte(id)), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
