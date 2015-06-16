package esxcloud

import (
	"encoding/json"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud/internal/rest"
	"io"
)

// Contains functionality for images API.
type ImagesAPI struct {
	client *Client
}

// Uploads a new image, reading from the specified image path.
// If options is nil, default options are used.
func (api *ImagesAPI) CreateFromFile(imagePath string, options *ImageCreateOptions) (task *Task, err error) {
	params := imageCreateOptionsToMap(options)
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
// If options is nil, default options are used.
func (api *ImagesAPI) Create(reader io.Reader, name string, options *ImageCreateOptions) (task *Task, err error) {
	params := imageCreateOptionsToMap(options)
	res, err := rest.MultipartUpload(api.client.httpClient, api.client.Endpoint+"/v1/images", reader, name, params)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := getTask(getError(res))
	return result, err
}

// Gets all images on this esxcloud instance.
func (api *ImagesAPI) GetAll() (images *Images, err error) {
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

// Gets details of image with the specified ID.
func (api *ImagesAPI) Get(imageID string) (image *Image, err error) {
	res, err := rest.Get(api.client.httpClient, api.client.Endpoint+"/v1/images/"+imageID)
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

// Deletes image with the specified ID.
func (api *ImagesAPI) Delete(imageID string) (task *Task, err error) {
	res, err := rest.Delete(api.client.httpClient, api.client.Endpoint+"/v1/images/"+imageID)
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
