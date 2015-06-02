package main

import (
	"errors"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
	"net/http"
	"os"
)

func CreateVM(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 6 {
		return nil, errors.New("Expected at least 6 arguments")
	}
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
	vmFlavor, ok := cloudProps["vm_flavor"].(string)
	if !ok {
		return nil, errors.New("Property 'vm_flavor' on cloud_properties is not a string or is not present")
	}
	diskFlavor, ok := cloudProps["disk_flavor"].(string)
	if !ok {
		return nil, errors.New("Property 'disk_flavor' on cloud_properties is not a string or is not present")
	}
	networks, ok := args[3].(map[string]interface{})
	if !ok {
		return nil, errors.New("Unexpected argument where networks should be")
	}
	// Ignore args[4] for now, which is disk_cids
	env, ok := args[5].(map[string]interface{})
	if !ok {
		return nil, errors.New("Unexpected argument where env should be")
	}

	ctx.Logger.Infof(
		"CreateVM with agent_id: '%v', stemcell_cid: '%v', cloud_properties: '%v', networks: '%v', env: '%v'",
		agentID, stemcellCID, cloudProps, networks, env)

	spec := &ec.VmCreateSpec{
		Name:          "bosh-vm",
		Flavor:        vmFlavor,
		SourceImageID: stemcellCID,
		AttachedDisks: []ec.AttachedDisk{
			ec.AttachedDisk{
				CapacityGB: 50, // Ignored
				Flavor:     diskFlavor,
				Kind:       "ephemeral-disk",
				Name:       "boot-disk",
				State:      "STARTED",
				BootDisk:   true,
			},
			ec.AttachedDisk{
				CapacityGB: 4, // Ignored
				Flavor:     diskFlavor,
				Kind:       "ephemeral-disk",
				Name:       "bosh-ephemeral-disk",
				State:      "STARTED",
				BootDisk:   false,
			},
		},
	}
	ctx.Logger.Infof("Creating VM with spec: %#v", spec)
	vmTask, err := ctx.Client.Projects.CreateVM(ctx.Config.ESXCloud.ProjectID, spec)
	if err != nil {
		return
	}
	ctx.Logger.Infof("Waiting on task: %#v", vmTask)
	vmTask, err = ctx.Client.Tasks.Wait(vmTask.ID)
	if err != nil {
		return
	}

	// Create and attach agent env ISO file
	envJson := createAgentEnv(ctx, agentID, vmTask.Entity.ID, spec.Name, networks, env)
	ctx.Logger.Infof("Creating agent env: %#v", envJson)
	isoPath, err := createEnvISO(envJson, ctx.Runner)
	if err != nil {
		return
	}
	defer os.Remove(isoPath)

	// Store env JSON as metadata so it can be picked up by attach/detach disk
	ctx.Logger.Info("Updating metadata for VM")
	err = putAgentEnvMetadata(vmTask.Entity.ID, envJson)
	if err != nil {
		return
	}

	ctx.Logger.Infof("Attaching ISO at path: %s", isoPath)
	attachTask, err := ctx.Client.VMs.AttachISO(vmTask.Entity.ID, isoPath)
	if err != nil {
		return
	}
	ctx.Logger.Infof("Waiting on task: %#v", attachTask)
	attachTask, err = ctx.Client.Tasks.Wait(attachTask.ID)
	if err != nil {
		return
	}

	op := &ec.VmOperation{Operation: "START_VM"}
	ctx.Logger.Info("Starting VM")
	onTask, err := ctx.Client.VMs.Operation(vmTask.Entity.ID, op)
	if err != nil {
		return
	}
	ctx.Logger.Infof("Waiting on task: %#v", onTask)
	onTask, err = ctx.Client.Tasks.Wait(onTask.ID)
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

	ctx.Logger.Infof("Deleting VM: %s", vmCID)

	ctx.Logger.Info("Detaching disks")
	// Detach any attached disks first
	disks, err := ctx.Client.Projects.FindDisks(ctx.Config.ESXCloud.ProjectID, nil)
	if err != nil {
		return
	}
	for _, disk := range disks.Items {
		for _, vmID := range disk.VMs {
			if vmID == vmCID {
				ctx.Logger.Infof("Detaching disk: %s", disk.ID)
				detachOp := &ec.VmDiskOperation{DiskID: disk.ID}
				detachTask, err := ctx.Client.VMs.DetachDisk(vmCID, detachOp)
				if err != nil {
					return nil, err
				}
				ctx.Logger.Infof("Waiting on task: %#v", detachTask)
				detachTask, err = ctx.Client.Tasks.Wait(detachTask.ID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	ctx.Logger.Info("Stopping VM")
	op := &ec.VmOperation{Operation: "STOP_VM"}
	offTask, err := ctx.Client.VMs.Operation(vmCID, op)
	if err != nil {
		return
	}
	ctx.Logger.Infof("Waiting on task: %#v", offTask)
	offTask, err = ctx.Client.Tasks.Wait(offTask.ID)
	if err != nil {
		return
	}

	ctx.Logger.Info("Deleting VM")
	task, err := ctx.Client.VMs.Delete(vmCID, true)
	if err != nil {
		return
	}
	ctx.Logger.Infof("Waiting on task: %#v", task)
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

	ctx.Logger.Infof("Determining if VM exists: %s", vmCID)
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

	ctx.Logger.Infof("Restarting VM: %s", vmCID)
	op := &ec.VmOperation{Operation: "RESTART_VM"}
	task, err := ctx.Client.VMs.Operation(vmCID, op)
	if err != nil {
		return
	}
	ctx.Logger.Infof("Waiting on task: %#v", task)
	_, err = ctx.Client.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return nil, nil
}
