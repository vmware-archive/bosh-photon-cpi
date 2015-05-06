package main

import (
	"errors"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
	"math"
)

func CreateDisk(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	size, ok := args[0].(int)
	if !ok {
		return nil, errors.New("Unexpected argument where size should be")
	}
	size = toGB(size)
	if size < 1 {
		return nil, errors.New("Must provide a size in MiB that rounds up to at least 1 GiB for esxcloud")
	}
	vmCID, ok := args[1].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where vm_cid should be")
	}

	diskSpec := &ec.DiskCreateSpec{
		Flavor:     ctx.Config.ESXCloud.DiskFlavor,
		Kind:       "persistent-disk",
		CapacityGB: size,
		Name:       "disk-for-vm-" + vmCID,
	}

	task, err := ctx.Client.Projects.CreateDisk(ctx.Config.ESXCloud.ProjectID, diskSpec)
	if err != nil {
		return
	}
	task, err = ctx.Client.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return task.Entity.ID, nil
}

func DeleteDisk(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	diskCID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where disk_cid should be")
	}
	task, err := ctx.Client.Disks.Delete(diskCID, true)
	if err != nil {
		return
	}
	task, err = ctx.Client.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return nil, nil
}

func toGB(mb int) int {
	return int(math.Ceil(float64(mb) / 1000.0))
}
