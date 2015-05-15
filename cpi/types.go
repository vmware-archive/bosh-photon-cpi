package cpi

import (
	"fmt"
	"github.com/esxcloud/bosh-esxcloud-cpi/cmd"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud"
)

type Context struct {
	Client *esxcloud.Client
	Config *Config
	Runner cmd.Runner
}

type Config struct {
	ESXCloud *ESXCloudConfig `json:"esxcloud"`
	Agent    *AgentConfig    `json:"agent"`
}

type AgentConfig struct {
	Mbus string   `json:"mbus"`
	NTP  []string `json:"ntp"`
}

type ESXCloudConfig struct {
	Target     string `json:"target"`
	ProjectID  string `json:"project"`
	TenantID   string `json:"tenant"`
	DiskFlavor string `json:"DiskFlavor"`
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
	Result interface{}    `json:"result,omitempty"`
	Error  *ResponseError `json:"error,omitempty"`
	Log    string         `json:"log,omitempty"`
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
	AgentID  string                 `json:"agent_id"`
	VM       VMSpec                 `json:"vm"`
	Mbus     string                 `json:"mbus"`
	NTP      []string               `json:"ntp"`
	Networks []interface{}          `json:"networks"`
	Env      map[string]interface{} `json:"env"`
}

type VMSpec struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
