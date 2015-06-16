package esxcloud

import (
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
)

// Contains functionality for status API.
type StatusAPI struct {
	client *Client
}

// Returns the status of an esxcloud endpoint.
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
