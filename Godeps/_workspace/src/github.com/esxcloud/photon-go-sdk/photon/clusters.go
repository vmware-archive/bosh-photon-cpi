package photon

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for clusters API.
type ClustersAPI struct {
	client *Client
}

var clusterUrl string = "/clusters/"

// Deletes a cluster with specified ID.
func (api *ClustersAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+clusterUrl+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets a cluster with the specified ID.
func (api *ClustersAPI) Get(id string) (cluster *Cluster, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+clusterUrl+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Cluster
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

// Gets vms for clusters with the specified ID
func (api *ClustersAPI) GetVMs(id string) (result *VMs, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+clusterUrl+id+"/vms", api.client.options.Token)
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

// Resize a cluster to specified count
func (api *ClustersAPI) Resize(id string, resize *ClusterResizeOperation) (task *Task, err error) {
	body, err := json.Marshal(resize)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+clusterUrl+id+"/resize", bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
