package photon

import (
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for auth API.
type AuthAPI struct {
	client *Client
}

var authUrl string = "/auth"

// Gets authentication info.
func (api *AuthAPI) Get() (info *AuthInfo, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+authUrl, api.client.options.Token)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	info = &AuthInfo{}
	err = json.NewDecoder(res.Body).Decode(info)
	return
}
