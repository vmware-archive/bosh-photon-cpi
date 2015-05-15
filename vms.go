package main

import (
	"errors"
	"fmt"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
	"net/http"
	"os"
)

func CreateVM(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 6 {
		return nil, errors.New("Expected at least 6 arguments")
	}
	fmt.Println(1)
	agentID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where agent_id should be")
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
	networks, ok := args[3].([]interface{})
	if !ok {
		return nil, errors.New("Unexpected argument where networks should be")
	}
	// Ignore args[4] for now, which is disk_cids
	env, ok := args[5].(map[string]interface{})
	if !ok {
		return nil, errors.New("Unexpected argument where env should be")
	}

	spec := &ec.VmCreateSpec{
		Name:          "bosh-vm",
		Flavor:        flavor,
		SourceImageID: stemcellCID,
	}
	vmTask, err := ctx.Client.Projects.CreateVM(ctx.Config.ESXCloud.ProjectID, spec)
	if err != nil {
		return
	}
	vmTask, err = ctx.Client.Tasks.Wait(vmTask.ID)
	if err != nil {
		return
	}

	// Create and attach agent env ISO file
	envJson := createAgentEnv(ctx, agentID, vmTask.Entity.ID, spec.Name, networks, env)
	isoPath, err := createEnvISO(envJson, ctx.Runner)
	if err != nil {
		return
	}
	defer os.Remove(isoPath)

	attachTask, err := ctx.Client.VMs.AttachISO(vmTask.Entity.ID, isoPath)
	if err != nil {
		return
	}
	attachTask, err = ctx.Client.Tasks.Wait(attachTask.ID)
	if err != nil {
		return
	}

	return vmTask.Entity.ID, nil
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
