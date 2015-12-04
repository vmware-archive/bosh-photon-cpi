package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/esxcloud/bosh-photon-cpi/cmd"
	"github.com/esxcloud/bosh-photon-cpi/cpi"
	"github.com/esxcloud/bosh-photon-cpi/logger"
	"github.com/esxcloud/photon-go-sdk/photon"
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
		res = createErrorResponse(cpi.NewCpiError(err, "Error reading from stdin"), "")
		return
	}

	req := &cpi.Request{}
	err = json.Unmarshal(reqBytes, req)
	if err != nil {
		res = createErrorResponse(cpi.NewCpiError(err, "Error deserializing JSON request from bosh"), "")
		return
	}

	configPath := flag.String("configPath", "", "Path to photon config file")
	flag.Parse()

	context, err := loadConfig(*configPath)
	if err != nil {
		res = createErrorResponse(cpi.NewCpiError(err, "Unable to load photon config from path '%s'", *configPath), "")
		return
	}

	// If there's an error with the logger, print it to stderr, but don't do anything
	// to prevent the CPI from running.
	if err != nil {
		os.Stderr.WriteString("Unable to create log file for photon CPI")
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
	clientConfig := &photon.ClientOptions{
		IgnoreCertificate: config.Photon.IgnoreCertificate,
		Token: config.Photon.Token,
	}
	ctx = &cpi.Context{
		Client: photon.NewClient(config.Photon.Target, clientConfig),
		Config: config,
		Runner: cmd.NewRunner(),
		Logger: logger.New(),
	}
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
			e := fmt.Errorf("%v", r)
			context.Logger.Error(e)
			result = createErrorResponse(e, context.Logger.LogData())
		}
	}()
	if fn, ok := actions[method]; ok {
		context.Logger.Infof("Begin action %s", method)
		context.Logger.Infof("Raw action arguments: %#v", args)

		res, err := fn(context, args)
		if err != nil {
			context.Logger.Errorf("Error encountered during action %s: %v", method, err)
			return createErrorResponse(err, context.Logger.LogData())
		}

		context.Logger.Infof("Action response: %#v", res)
		context.Logger.Infof("End action %s", method)
		return createResponse(res, context.Logger.LogData())
	} else {
		e := cpi.NewBoshError(cpi.NotSupportedError, false, "Method %s not supported in photon CPI.", method)
		context.Logger.Error(e)
		return createErrorResponse(e, context.Logger.LogData())
	}
	return
}

func createResponse(result interface{}, logData string) []byte {
	res := &cpi.Response{Result: result, Log: logData, Error: nil}
	resBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return resBytes
}

func createErrorResponse(err error, logData string) []byte {
	res := &cpi.Response{
		Error: &cpi.ResponseError{
			Message: err.Error(),
		},
		Log: logData,
	}

	switch t := err.(type) {
	// If caller throws BoshError specifically, respect type and canRetry from caller
	case cpi.BoshError:
		res.Error.Type = t.Type()
		res.Error.CanRetry = t.CanRetry()
	// An API error or a task in error state cannot be retried
	case photon.ApiError, photon.TaskError:
		res.Error.Type = cpi.CloudError
		res.Error.CanRetry = false
	// Task timeout errors and unknown HTTP errors can likely be retried
	case photon.HttpError, photon.TaskTimeoutError:
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
