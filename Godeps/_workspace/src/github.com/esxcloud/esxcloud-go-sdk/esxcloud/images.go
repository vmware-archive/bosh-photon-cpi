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

type ImageCreateOptions struct {
	ReplicationType string
}

type Images struct {
	Items []Image `json:"items"`
}

// Uploads a new image, reading from the specified image path.
// If opts is nil, default options are used.
func (api *ImagesAPI) CreateFromFile(imagePath string, opts *ImageCreateOptions) (task *Task, err error) {
	params := imageCreateOptionsToMap(opts)
	res, err := rest.MultipartUploadFile(api.client.httpClient, api.client.Endpoint+"/v1/images", imagePath, params)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := getTask(getError(res))
	return result, err
}

// Uploads a new image, reading from the specified io.Reader.
// Name is a descriptive name of the image, it is used in the filename field of the Content-Disposition header,
// and does not need to be unique.
// If opts is nil, default options are used.
func (api *ImagesAPI) Create(reader io.Reader, name string, opts *ImageCreateOptions) (task *Task, err error) {
	params := imageCreateOptionsToMap(opts)
	res, err := rest.MultipartUpload(api.client.httpClient, api.client.Endpoint+"/v1/images", reader, name, params)
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

func defaultImageCreateOptions() *ImageCreateOptions {
	return &ImageCreateOptions{ReplicationType: "EAGER"}
}

func imageCreateOptionsToMap(opts *ImageCreateOptions) map[string]string {
	if opts == nil {
		return imageCreateOptionsToMap(defaultImageCreateOptions())
	}
	return map[string]string{
		"ImageReplication": opts.ReplicationType,
	}
}
