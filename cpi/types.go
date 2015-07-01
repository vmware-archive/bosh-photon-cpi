package cpi

import (
	"fmt"
	"github.com/esxcloud/bosh-esxcloud-cpi/cmd"
	"github.com/esxcloud/bosh-esxcloud-cpi/logger"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud"
)

type Context struct {
	Client *esxcloud.Client
	Config *Config
	Runner cmd.Runner
	Logger logger.Logger
}

type Config struct {
	ESXCloud *ESXCloudConfig `json:"esxcloud"`
	Agent    *AgentConfig    `json:"agent"`
}

type AgentConfig struct {
	Mbus      string        `json:"mbus"`
	NTP       []string      `json:"ntp"`
	Blobstore BlobstoreSpec `json:"blobstore"`
}

type ESXCloudConfig struct {
	Target            string `json:"target"`
	ProjectID         string `json:"project"`
	TenantID          string `json:"tenant"`
	IgnoreCertificate bool   `json:"ignore_cert"`
}

type ActionFn func(*Context, []interface{}) (interface{}, error)

type BoshErrorType string

const (
	CloudError          BoshErrorType = "Bosh::Clouds::CloudError"
	CpiError            BoshErrorType = "Bosh::Clouds::CpiError"
	NotImplementedError BoshErrorType = "Bosh::Clouds::NotImplemented"
)

type Request struct {
	Method    string        `json:"method"`
	Arguments []interface{} `json:"arguments"`
}

type Response struct {
	Result interface{}    `json:"result"`
	Error  *ResponseError `json:"error"`
	Log    string         `json:"log"`
}

type ResponseError struct {
	Type     BoshErrorType `json:"type"`
	Message  string        `json:"message"`
	CanRetry bool          `json:"ok_to_retry"`
}

type BoshError interface {
	Type() BoshErrorType
	CanRetry() bool
}

type boshError struct {
	errorType BoshErrorType
	canRetry  bool
	message   string
}

func (e boshError) Type() BoshErrorType {
	return e.errorType
}

func (e boshError) CanRetry() bool {
	return e.canRetry
}

func (e boshError) Error() string {
	return e.message
}

func NewBoshError(errorType BoshErrorType, canRetry bool, format string, args ...interface{}) error {
	return &boshError{errorType, canRetry, fmt.Sprintf(format, args...)}
}

func NewCpiError(cause interface{}, format string, args ...interface{}) error {
	return &boshError{CpiError, false, fmt.Sprintf("CPI error: '%s' | Caused by: '%v'", fmt.Sprintf(format, args...), cause)}
}

type Network struct {
	Type            string                 `json:"type"`
	IP              string                 `json:"ip"`
	Netmask         string                 `json:"netmask"`
	Gateway         string                 `json:"gateway"`
	DNS             []string               `json:"dns"`
	Default         []string               `json:"default"`
	MAC             string                 `json:"mac"`
	CloudProperties map[string]interface{} `json:"cloud_properties"`
}

type AgentEnv struct {
	AgentID   string                 `json:"agent_id"`
	VM        VMSpec                 `json:"vm"`
	Mbus      string                 `json:"mbus"`
	NTP       []string               `json:"ntp"`
	Networks  map[string]interface{} `json:"networks"`
	Env       map[string]interface{} `json:"env"`
	Disks     map[string]interface{} `json:"disks"`
	Blobstore BlobstoreSpec          `json:"blobstore"`
}

type VMSpec struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type BlobstoreSpec struct {
	Provider string                 `json:"provider"`
	Options  map[string]interface{} `json:"options"`
}
