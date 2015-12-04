package photon

import (
	"encoding/json"

	"github.com/esxcloud/photon-go-sdk/photon/internal/rest"
)

// Contains functionality for hosts API.
type ResourceTicketsAPI struct {
	client *Client
}

var resourceTicketUrl string = "/resource-tickets/"

// Gets all tasks with the specified resource ticket ID, using options to filter the results.
// If options is nil, no filtering will occur.
func (api *ResourceTicketsAPI) GetTasks(id string, options *TaskGetOptions) (result *TaskList, err error) {
	uri := api.client.Endpoint + resourceTicketUrl + id + "/tasks"
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
