package esxcloud

import (
	"bytes"
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

type FlavorsAPI struct {
	client *Client
}

// Options used for find/get APIs
type FlavorOptions struct {
	Name string
	Kind string
}

type FlavorCreateSpec struct {
	Cost []QuotaLineItem `json:"cost"`
	Kind string          `json:"kind"`
	Name string          `json:"name"`
}

type Flavor struct {
	Cost     []QuotaLineItem `json:"cost"`
	Kind     string          `json:"kind"`
	Name     string          `json:"name"`
	ID       string          `json:"id"`
	Tags     []string        `json:"tags"`
	SelfLink string          `json:"selfLink"`
}

type FlavorList struct {
	Items []Flavor `json:"items"`
}

func (api *FlavorsAPI) Create(spec *FlavorCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+"/v1/flavors", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

func (api *FlavorsAPI) Get(id string) (flavor *Flavor, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/flavors/"+id)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	flavor = &Flavor{}
	err = json.NewDecoder(res.Body).Decode(flavor)
	return
}

// Gets flavors using options to filter results. Returns all flavors if options is nil.
func (api *FlavorsAPI) Find(options *FlavorOptions) (flavors *FlavorList, err error) {
	uri := api.client.Endpoint + "/v1/flavors"
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
	flavors = &FlavorList{}
	err = json.NewDecoder(res.Body).Decode(flavors)
	return
}

func (api *FlavorsAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/flavors/"+id)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}
