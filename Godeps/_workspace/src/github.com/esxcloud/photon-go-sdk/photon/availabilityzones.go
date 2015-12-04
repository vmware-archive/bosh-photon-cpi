package photon

import (
	"bytes"
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for availability zones API.
type AvailabilityZonesAPI struct {
	client *Client
}

var availabilityzoneUrl string = "/availabilityzones"

// Creates availability zone.
func (api *AvailabilityZonesAPI) Create(availabilityzoneSpec *AvailabilityZoneCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(availabilityzoneSpec)
	if err != nil {
		return
	}
	res, err := rest.Post(api.client.httpClient, api.client.Endpoint+availabilityzoneUrl, bytes.NewReader(body), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets availability zone with the specified ID.
func (api *AvailabilityZonesAPI) Get(id string) (availabilityzone *AvailabilityZone, err error) {
	res, err := rest.Get(api.client.httpClient, api.getEntityUrl(id), api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	availabilityzone = &AvailabilityZone{}
	err = json.NewDecoder(res.Body).Decode(availabilityzone)
	return
}

// Returns all availability zones on an photon instance.
func (api *AvailabilityZonesAPI) GetAll() (result *AvailabilityZones, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+availabilityzoneUrl, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &AvailabilityZones{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Deletes the availability zone with specified ID.
func (api *AvailabilityZonesAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+availabilityzoneUrl+"/"+id, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets all tasks with the specified availability zone ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *AvailabilityZonesAPI) GetTasks(id string, options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint + availabilityzoneUrl + "/" + id + "/tasks"
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

func (api *AvailabilityZonesAPI) getEntityUrl(id string) (url string) {
	return api.client.Endpoint + availabilityzoneUrl + "/" + id
}
