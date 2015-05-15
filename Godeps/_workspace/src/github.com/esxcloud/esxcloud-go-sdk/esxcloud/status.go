package esxcloud

import (
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

type StatusAPI struct {
	client *Client
}

type Component struct {
	Component string
	Message   string
	Status    string
}

type Status struct {
	Status     string
	Components []Component
}

// Returns the status of an esxcloud endpoint. Returns ApiError or HttpError
// in the event of an API error or unknown HTTP error, respectively.
func (api *StatusAPI) Get() (status *Status, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/status")
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	status = &Status{}
	err = json.NewDecoder(res.Body).Decode(status)
	return
}
