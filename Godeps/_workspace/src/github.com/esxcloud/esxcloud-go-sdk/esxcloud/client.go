package esxcloud

import (
	"net/http"
	"time"
)

// Represents stateless context needed to call esxcloud APIs.
type Client struct {
	options    ClientOptions
	httpClient *http.Client
	Endpoint   string
	Status     *StatusAPI
	Tenants    *TenantsAPI
	Tasks      *TasksAPI
	Projects   *ProjectsAPI
	Flavors    *FlavorsAPI
	Images     *ImagesAPI
	Disks      *DisksAPI
	VMs        *VmAPI
}

// Options for Client
type ClientOptions struct {
	// When using the Tasks.Wait APIs, defines the duration of how long
	// the SDK should continue to poll the server. Default is 30 minutes.
	// TasksAPI.WaitTimeout() can be used to specify timeout on
	// individual calls.
	TaskPollTimeout time.Duration

	// For tasks APIs, defines the delay between each polling attempt.
	// Default is 100 milliseconds.
	taskPollDelay time.Duration

	// For tasks APIs, defines the number of retries to make in the event
	// of an error. Default is 3.
	taskRetryCount int
}

// Creates a new ESXCloud client with specified options. If options
// is nil, default options will be used.
func NewClient(endpoint string, options *ClientOptions) (c *Client) {
	c = &Client{Endpoint: endpoint, httpClient: &http.Client{}}
	c.Status = &StatusAPI{c}
	c.Tenants = &TenantsAPI{c}
	c.Tasks = &TasksAPI{c}
	c.Projects = &ProjectsAPI{c}
	c.Flavors = &FlavorsAPI{c}
	c.Images = &ImagesAPI{c}
	c.Disks = &DisksAPI{c}
	c.VMs = &VmAPI{c}

	if options == nil {
		options = &ClientOptions{
			TaskPollTimeout: 30 * time.Minute,
			taskPollDelay:   100 * time.Millisecond,
			taskRetryCount:  3,
		}
	}
	// Ensure a copy of options is made, rather than using a pointer
	// which may change out from underneath if misused by the caller.
	c.options = *options
	return
}

// Creates a new ESXCloud client with specified options and http.Client.
// Useful for functional testing where http calls must be mocked out.
// If options is nil, default options will be used.
func NewTestClient(endpoint string, options *ClientOptions, httpClient *http.Client) (c *Client) {
	c = NewClient(endpoint, options)
	c.httpClient = httpClient
	return
}
