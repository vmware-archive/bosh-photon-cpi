package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/esxcloud/bosh-esxcloud-cpi/cmd"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	"github.com/esxcloud/esxcloud-go-sdk/esxcloud"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func main() {
	actions := map[string]cpi.ActionFn{
		"create_stemcell": CreateStemcell,
		"delete_stemcell": DeleteStemcell,
		"create_disk":     CreateDisk,
		"delete_disk":     DeleteDisk,
		"has_disk":        HasDisk,
		"attach_disk":     AttachDisk,
		"detach_disk":     DetachDisk,
		"create_vm":       CreateVM,
		"delete_vm":       DeleteVM,
		"has_vm":          HasVM,
		"restart_vm":      RestartVM,
		"set_vm_metadata": SetVmMetadata,
	}

	var res []byte
	defer func() { os.Stdout.Write(res) }()

	reqBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		res = createErrorResponse(cpi.NewCpiError(err, "Error reading from stdin"))
		return
	}

	req := &cpi.Request{}
	err = json.Unmarshal(reqBytes, req)
	if err != nil {
		res = createErrorResponse(cpi.NewCpiError(err, "Error deserializing JSON request from bosh"))
		return
	}

	configPath := flag.String("configPath", "", "Path to esxcloud config file")
	flag.Parse()

	context, err := loadConfig(*configPath)
	if err != nil {
		res = createErrorResponse(cpi.NewCpiError(err, "Unable to load esxcloud config from path '%s'", *configPath))
		return
	}

	res = dispatch(context, actions, strings.ToLower(req.Method), req.Arguments)
}

func loadConfig(filePath string) (ctx *cpi.Context, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	config := &cpi.Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return
	}
	ctx = &cpi.Context{Client: esxcloud.NewClient(config.ESXCloud.Target, nil), Config: config, Runner: cmd.NewRunner()}
	return
}

func dispatch(context *cpi.Context, actions map[string]cpi.ActionFn, method string, args []interface{}) (result []byte) {
	// Attempt to recover from any panic that may occur during API calls
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				// Don't even try to recover severe runtime errors
				panic(r)
			}
			result = createErrorResponse(errors.New(fmt.Sprintf("%v", r)))
		}
	}()
	if fn, ok := actions[method]; ok {
		res, err := fn(context, args)
		if err != nil {
			return createErrorResponse(err)
		}
		return createResponse(res)
	} else {
		return createErrorResponse(
			cpi.NewBoshError(cpi.NotImplementedError, false, "Method %s not implemented in esxcloud CPI.", method))
	}
	return
}

func createResponse(result interface{}) []byte {
	res := &cpi.Response{Result: result}
	resBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return resBytes
}

func createErrorResponse(err error) []byte {
	res := &cpi.Response{
		Error: &cpi.ResponseError{
			Message: err.Error(),
		}}

	switch t := err.(type) {
	// If caller throws BoshError specifically, respect type and canRetry from caller
	case cpi.BoshError:
		res.Error.Type = t.Type()
		res.Error.CanRetry = t.CanRetry()
	// An API error or a task in error state cannot be retried
	case esxcloud.ApiError, esxcloud.TaskError:
		res.Error.Type = cpi.CloudError
		res.Error.CanRetry = false
	// Task timeout errors and unknown HTTP errors can likely be retried
	case esxcloud.HttpError, esxcloud.TaskTimeoutError:
		res.Error.Type = cpi.CloudError
		res.Error.CanRetry = true
	// Assume unknown errors are CPI errors that cannnot be retried
	default:
		res.Error.Type = cpi.CpiError
		res.Error.CanRetry = false
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return resBytes
}
