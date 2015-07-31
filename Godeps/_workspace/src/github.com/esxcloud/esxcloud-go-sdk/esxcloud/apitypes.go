package esxcloud

import (
	"fmt"
)

type Entity struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

// Represents an error from the esxcloud API.
type ApiError struct {
	Code           string                 `json:"code"`
	Data           map[string]interface{} `json:"data"`
	Message        string                 `json:"message"`
	HttpStatusCode int                    `json:"-"` // Not part of API contract
}

// Implement Go error interface for ApiError
func (e ApiError) Error() string {
	return fmt.Sprintf(
		"esxcloud: { HTTP status: '%v', code: '%v', message: '%v', data: '%v' }",
		e.HttpStatusCode,
		e.Code,
		e.Message,
		e.Data)
}

// Used to represent a generic HTTP error, i.e. an unexpected HTTP 500.
type HttpError struct {
	StatusCode int
	Message    string
}

// Implementation of error interface for HttpError
func (e HttpError) Error() string {
	return fmt.Sprintf("esxcloud: HTTP %d: %v", e.StatusCode, e.Message)
}

// Represents an ESXCloud task that has entered into an error state.
// ESXCloud task errors can be caught and type-checked against with
// the usual Go idiom.
type TaskError struct {
	ID string
}

// Implement Go error interface for TaskError.
func (e TaskError) Error() string {
	return fmt.Sprintf("esxcloud: Task '%s' is in error state. "+
		"Examine task for full details.", e.ID)
}

// An error representing a timeout while waiting for a task to complete.
type TaskTimeoutError struct {
	ID string
}

// Implement Go error interface for TaskTimeoutError.
func (e TaskTimeoutError) Error() string {
	return fmt.Sprintf("esxcloud: Timed out waiting for task '%s'. "+
		"Task may not be in error state, examine task for full details.", e.ID)
}

// Represents an operation (Step) within a Task.
type Step struct {
	ID                 string                 `json:"id"`
	Operation          string                 `json:"operation,omitempty"`
	State              string                 `json:"state"`
	StartedTime        int64                  `json:"startedTime"`
	EndTime            int64                  `json:"endTime,omitempty"`
	QueuedTime         int64                  `json:"queuedTime"`
	Sequence           int                    `json:"sequence,omitempty"`
	ResourceProperties map[string]interface{} `json:"resourceProperties,omitempty"`
	Errors             []ApiError             `json:"errors,omitempty"`
	Options            map[string]interface{} `json:"options,omitempty"`
	SelfLink           string                 `json:"selfLink,omitempty"`
}

// Represents an asynchronous task.
type Task struct {
	ID          string `json:"id"`
	Operation   string `json:"operation,omitempty"`
	State       string `json:"state"`
	StartedTime int64  `json:"startedTime"`
	EndTime     int64  `json:"endTime,omitempty"`
	QueuedTime  int64  `json:"queuedTime"`
	Entity      Entity `json:"entity,omitempty"`
	SelfLink    string `json:"selfLink,omitempty"`
	Steps       []Step `json:"steps,omitempty"`
}

// Represents multiple tasks returned by the API.
type TaskList struct {
	Items []Task `json:"items"`
}

// Options for GetTasks API.
type TaskGetOptions struct {
	State string
	Kind  string
}

type BaseCompact struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type QuotaLineItem struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
	Key   string  `json:"key"`
}

// Creation spec for locality.
type LocalitySpec struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

// Creation spec for disks.
type DiskCreateSpec struct {
	Flavor     string         `json:"flavor"`
	Kind       string         `json:"kind"`
	CapacityGB int            `json:"capacityGb"`
	Affinities []LocalitySpec `json:"localitySpec,omitempty"`
	Name       string         `json:"name"`
	Tags       []string       `json:"tags,omitempty"`
}

// Represents a persistent disk.
type PersistentDisk struct {
	Flavor     string          `json:"flavor"`
	Cost       []QuotaLineItem `json:"cost"`
	Kind       string          `json:"kind"`
	Datastore  string          `json:"datastore,omitempty"`
	CapacityGB int             `json:"capacityGb,omitempty"`
	Name       string          `json:"name"`
	State      string          `json:"state"`
	ID         string          `json:"id"`
	VMs        []string        `json:"vms"`
	Tags       []string        `json:"tags,omitempty"`
	SelfLink   string          `json:"selfLink,omitempty"`
}

// Represents multiple persistent disks returned by the API.
type DiskList struct {
	Items []PersistentDisk `json:"items"`
}

// Creation spec for projects.
type ProjectCreateSpec struct {
	ResourceTicket ResourceTicketReservation `json:"resourceTicket"`
	Name           string                    `json:"name"`
}

// Represents multiple projects returned by the API.
type ProjectList struct {
	Items []ProjectCompact `json:"items"`
}

// Compact representation of projects.
type ProjectCompact struct {
	Kind           string        `json:"kind"`
	ResourceTicket ProjectTicket `json:"resourceTicket"`
	Name           string        `json:"name"`
	ID             string        `json:"id"`
	Tags           []string      `json:"tags"`
	SelfLink       string        `json:"selfLink"`
}

type ProjectTicket struct {
	TenantTicketID   string          `json:"tenantTicketId"`
	Usage            []QuotaLineItem `json:"usage"`
	TenantTicketName string          `json:"tenantTicketName"`
	Limits           []QuotaLineItem `json:"limits"`
}

// Represents an image.
type Image struct {
	Size            int64          `json:"size"`
	Kind            string         `json:"kind"`
	Name            string         `json:"name"`
	State           string         `json:"state"`
	ID              string         `json:"id"`
	Tags            []string       `json:"tags"`
	SelfLink        string         `json:"selfLink"`
	Settings        []ImageSetting `json:"settings"`
	ReplicationType string         `json:"replicationType"`
}

// Represents an image setting
type ImageSetting struct {
	Name         string `json:"name"`
	DefaultValue string `json:"defaultValue"`
}

// Creation spec for images.
type ImageCreateOptions struct {
	ReplicationType string
}

// Represents multiple images returned by the API.
type Images struct {
	Items []Image `json:"items"`
}

// Represents a component with status.
type Component struct {
	Component string
	Message   string
	Status    string
}

// Represents status of the esxcloud system.
type Status struct {
	Status     string
	Components []Component
}

// Represents a single tenant.
type Tenant struct {
	Projects        []BaseCompact `json:"projects"`
	ResourceTickets []BaseCompact `json:"resourceTickets"`
	Kind            string        `json:"kind"`
	Name            string        `json:"name"`
	ID              string        `json:"id"`
	SelfLink        string        `json:"selfLink"`
	Tags            []string      `json:"tags"`
}

// Represents multiple tenants returned by the API.
type Tenants struct {
	Items []Tenant `json:"items"`
}

// Creation spec for tenants.
type TenantCreateSpec struct {
	Name string `json:"name"`
}

// Creation spec for resource tickets.
type ResourceTicketCreateSpec struct {
	Name   string          `json:"name"`
	Limits []QuotaLineItem `json:"limits"`
}

// Represents a single resource ticket.
type ResourceTicket struct {
	Kind     string          `json:"kind"`
	Usage    []QuotaLineItem `json:"usage"`
	TenantId string          `json:"tenantId"`
	Name     string          `json:"name"`
	ID       string          `json:"id"`
	Limits   []QuotaLineItem `json:"limits"`
	Tags     []string        `json:"tags"`
	SelfLink string          `json:"selfLink"`
}

// Represents multiple resource tickets returned by the API.
type ResourceList struct {
	Items []ResourceTicket `json:"items"`
}

// Represents a resource reservation on a resource ticket.
type ResourceTicketReservation struct {
	Name   string          `json:"name"`
	Limits []QuotaLineItem `json:"limits"`
}

// Creation spec for VMs.
type VmCreateSpec struct {
	Flavor        string         `json:"flavor"`
	SourceImageID string         `json:"sourceImageId"`
	AttachedDisks []AttachedDisk `json:"attachedDisks"`
	Affinities    []LocalitySpec `json:"affinities,omitempty"`
	Name          string         `json:"name"`
	Tags          []string       `json:"tags,omitempty"`
}

// Represents possible operations for VMs. Valid values include:
// START_VM, STOP_VM, RESTART_VM, SUSPEND_VM, RESUME_VM
type VmOperation struct {
	Operation string                 `json:"operation"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// Represents metadata that can be set on a VM.
type VmMetadata struct {
	Metadata map[string]string `json:"metadata"`
}

// Represents a single attached disk.
type AttachedDisk struct {
	Flavor     string `json:"flavor"`
	Kind       string `json:"kind"`
	CapacityGB int    `json:"capacityGb"`
	Name       string `json:"name"`
	State      string `json:"state"`
	ID         string `json:"id,omitempty"`
	BootDisk   bool   `json:"bootDisk"`
}

// Represents a single VM.
type VM struct {
	SourceImageID string            `json:"sourceImageId,omitempty"`
	Cost          []QuotaLineItem   `json:"cost"`
	Kind          string            `json:"kind"`
	AttachedDisks []AttachedDisk    `json:"attachedDisks"`
	Datastore     string            `json:"datastore,omitempty"`
	AttachedISOs  []ISO             `json:"attachedIsos,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	SelfLink      string            `json:"selfLink,omitempty"`
	Flavor        string            `json:"flavor"`
	Host          string            `json:"host,omitempty"`
	Name          string            `json:"name"`
	State         string            `json:"string"`
	ID            string            `json:"id"`
}

// Represents multiple VMs returned by the API.
type VMs struct {
	Items []VM `json:"items"`
}

// Represents an ISO.
type ISO struct {
	Size int64  `json:"size,omitempty"`
	Kind string `json:"kind,omitempty"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

// Represents operations for disks.
type VmDiskOperation struct {
	DiskID    string                 `json:"diskId"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// Creation spec for flavors.
type FlavorCreateSpec struct {
	Cost []QuotaLineItem `json:"cost"`
	Kind string          `json:"kind"`
	Name string          `json:"name"`
}

// Represents a single flavor.
type Flavor struct {
	Cost     []QuotaLineItem `json:"cost"`
	Kind     string          `json:"kind"`
	Name     string          `json:"name"`
	ID       string          `json:"id"`
	Tags     []string        `json:"tags"`
	SelfLink string          `json:"selfLink"`
}

// Represents multiple flavors returned by the API.
type FlavorList struct {
	Items []Flavor `json:"items"`
}

// Creation spec for hosts.
type HostCreateSpec struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Metadata interface{} `json:"metadata,omitempty"`
	Address  string      `json:"address"`
	Tags     []string    `json:"usageTags"`
}

// Represents a host
type Host struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Address  string      `json:"address"`
	Kind     string      `json:"kind"`
	ID       string      `json:"id"`
	Tags     []string    `json:"usageTags"`
	Metadata interface{} `json:"metadata,omitempty"`
	SelfLink string      `json:"selfLink"`
	State    string      `json:"state"`
}

// Represents multiple hosts returned by the API.
type Hosts struct {
	Items []Host `json:"items"`
}

// Creation spec for deployments.
type DeploymentCreateSpec struct {
	NTPEndpoint             interface{} `json:"ntpEndpoint"`
	UseImageDatastoreForVMs bool        `json:"useImageDatastoreForVMs"`
	SyslogEndpoint          interface{} `json:"syslogEndpoint"`
	ImageDatastore          string      `json:"imageDatastore"`
	Auth                    *AuthInfo   `json:"auth"`
}

// Represents a deployment
type Deployment struct {
	NTPEndpoint             string    `json:"ntpEndpoint,omitempty"`
	UseImageDatastoreForVMs bool      `json:"useImageDatastoreForVMs,omitempty"`
	Auth                    *AuthInfo `json:"auth"`
	Kind                    string    `json:"kind"`
	SyslogEndpoint          string    `json:"syslogEndpoint,omitempty"`
	State                   string    `json:"state"`
	ID                      string    `json:"id"`
	ImageDatastore          string    `json:"imageDatastore"`
	SelfLink                string    `json:"selfLink"`
}

// Represents multiple deployments returned by the API.
type Deployments struct {
	Items []Deployment `json:"items"`
}

// Represents authentication information
type AuthInfo struct {
	Password string `json:"password,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	Tenant   string `json:"tenant,omitempty"`
	Enabled  bool   `json:"enabled"`
	Username string `json:"username,omitempty"`
}
