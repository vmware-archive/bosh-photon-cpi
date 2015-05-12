package main

import (
	"errors"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
)

func CreateStemcell(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("Expected at least 1 argument")
	}
	imagePath, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where image_path should be")
	}
	task, err := ctx.Client.Images.CreateFromFile(imagePath)
	if err != nil {
		return
	}
	task, err = ctx.Client.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return task.Entity.ID, nil
}

func DeleteStemcell(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("Expected at least 1 argument")
	}
	stemcellCID, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where stemcell_cid should be")
	}
	task, err := ctx.Client.Images.Delete(stemcellCID)
	if err != nil {
		return
	}
	task, err = ctx.Client.Tasks.Wait(task.ID)
	if err != nil {
		return
	}
	return nil, nil
}
