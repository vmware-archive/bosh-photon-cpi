package main

import (
	"errors"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
	"net/http"
)

func CreateVM(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 3 {
		return nil, errors.New("Expected at least 3 arguments")
	}
	stemcellCID, ok := args[1].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where stemcell_cid should be")
	}
	cloudProps, ok := args[2].(map[string]interface{})
	if !ok {
		return nil, errors.New("Unexpected argument where cloud_properties should be")
	}
	flavor, ok := cloudProps["flavor"].(string)
	if !ok {
		return nil, errors.New("Property 'flavor' on cloud_properties is not a string")
	}

	spec := &ec.VmCreateSpec{
		Flavor:        flavor,
		SourceImageID: stemcellCID,
	}
	task, err := ctx.Client.Projects.CreateVM(ctx.Config.ESXCloud.ProjectID, spec)
	if err != nil {
		return
	}
	task, err = ctx.Client.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return task.Entity.ID, nil
}

func DeleteVM(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("Expected at least 1 argument")
	}
	vmCID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where vm_cid should be")
	}
	task, err := ctx.Client.VMs.Delete(vmCID, true)
	if err != nil {
		return
	}
	_, err = ctx.Client.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return nil, nil
}

func HasVM(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("Expected at least 1 argument")
	}
	vmCID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where vm_cid should be")
	}
	_, err = ctx.Client.VMs.Get(vmCID)
	if err != nil {
		apiErr, ok := err.(ec.ApiError)
		if ok && apiErr.HttpStatusCode == http.StatusNotFound {
			return false, nil
		}
		return nil, err
	}
	return true, nil
}

func RestartVM(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("Expected at least 1 argument")
	}
	vmCID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where vm_cid should be")
	}
	op := &ec.VmOperation{Operation: "RESTART_VM"}
	task, err := ctx.Client.VMs.Operation(vmCID, op)
	if err != nil {
		return
	}
	_, err = ctx.Client.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return nil, nil
}
