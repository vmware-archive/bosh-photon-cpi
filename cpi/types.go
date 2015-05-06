package cpi

import (
	"fmt"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud"
)

type Context struct {
	Client *esxcloud.Client
	Config *Config
}

type Config struct {
	ESXCloud *ESXCloudConfig `json:"ESXCloud"`
}

type ESXCloudConfig struct {
	APIFE      string `json:"APIFE"`
	ProjectID  string `json:"ProjectID"`
	TenantID   string `json:"TenantID"`
	DiskFlavor string `json:"DiskFlavor"`
	VMFlavor   string `json:"VMFlavor"`
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
