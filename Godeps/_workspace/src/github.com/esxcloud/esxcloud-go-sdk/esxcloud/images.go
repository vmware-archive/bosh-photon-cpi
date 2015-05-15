package esxcloud

import (
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
	"io"
)

type ImagesAPI struct {
	client *Client
}

type Image struct {
	Size     int64    `json:"size"`
	Kind     string   `json:"kind"`
	Name     string   `json:"name"`
	State    string   `json:"state"`
	ID       string   `json:"id"`
	Tags     []string `json:"tags"`
	SelfLink string   `json:"selfLink"`
}

type Images struct {
	Items []Image `json:"items"`
}

func (api *ImagesAPI) CreateFromFile(imagePath string) (task *Task, err error) {
	res, err := rest.MultipartUploadFile(api.client.httpClient, api.client.Endpoint+"/v1/images", imagePath)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := getTask(getError(res))
	return result, err
}

func (api *ImagesAPI) Create(reader io.Reader) (task *Task, err error) {
	res, err := rest.MultipartUpload(api.client.httpClient, api.client.Endpoint+"/v1/images", reader)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := getTask(getError(res))
	return result, err
}

func (api *ImagesAPI) GetAll(client Client) (images *Images, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/images")
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Images
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

func (api *ImagesAPI) Get(id string) (image *Image, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/images/"+id)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Image
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

func (api *ImagesAPI) Delete(id string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/images/"+id)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := getTask(getError(res))
	return result, err
}
