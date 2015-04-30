package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/esxcloud/bosh-esxcloud-cpi/stemcell"
	. "github.com/esxcloud/bosh-esxcloud-cpi/types"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type actionFn func([]interface{}) (interface{}, error)

func main() {
	actions := map[string]actionFn{
		"create_stemcell": stemcell.Create,
		"delete_stemcell": stemcell.Delete,
	}

	reqBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic("Error reading from stdin")
	}

	req := &Request{}
	err = json.Unmarshal(reqBytes, req)
	if err != nil {
		panic("Error deserializing JSON request from bosh")
	}

	res := dispatch(actions, strings.ToLower(req.Method), req.Arguments)
	os.Stdout.Write(res)
}

func dispatch(actions map[string]actionFn, method string, args []interface{}) (result []byte) {
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
		res, err := fn(args)
		if err != nil {
			return createErrorResponse(err)
		}
		return createResponse(res)
	} else {
		return createErrorResponse(
			NewBoshError(NotImplementedError, false, "Method %s not implemented in esxcloud CPI", method))
	}
	return
}

func createResponse(result interface{}) []byte {
	res := &Response{Result: result}
	resBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return resBytes
}

func createErrorResponse(err error) []byte {
	res := &Response{
		Error: &ResponseError{
			Message: err.Error(),
		}}

	if typedErr, ok := err.(BoshError); ok {
		res.Error.Type = typedErr.Type()
		res.Error.CanRetry = typedErr.CanRetry()
	} else {
		res.Error.Type = CloudError
		res.Error.CanRetry = false
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return resBytes
}
