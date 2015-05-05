package main

import (
	"errors"
	. "github.com/esxcloud/bosh-esxcloud-cpi/types"
)

func CreateStemcell(ctx *Context, args []interface{}) (result interface{}, err error) {
	imagePath, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where image_path should be")
	}
	task, err := ctx.ECClient.Images.Create(imagePath)
	if err != nil {
		return
	}
	task, err = ctx.ECClient.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return task.Entity.ID, nil
}

func DeleteStemcell(ctx *Context, args []interface{}) (result interface{}, err error) {
	stemcellCID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where stemcell_cid should be")
	}
	task, err := ctx.ECClient.Images.Delete(stemcellCID)
	if err != nil {
		return
	}
	task, err = ctx.ECClient.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return nil, nil
}
