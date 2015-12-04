package photon

import (
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for status API.
type StatusAPI struct {
	client *Client
}

var statusUrl string = "/status"

// Returns the status of an photon endpoint.
func (api *StatusAPI) Get() (status *Status, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+statusUrl, api.client.options.Token)
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
