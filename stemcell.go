package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	"os"
	"path/filepath"
)

func CreateStemcell(ctx *cpi.Context, args []interface{}) (result interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("Expected at least 1 argument")
	}
	imagePath, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where image_path should be")
	}
	stemcell, err := newStemcell(imagePath)
	if err != nil {
		return
	}
	defer stemcell.Close()
	task, err := ctx.Client.Images.Create(stemcell, filepath.Base(imagePath), nil)
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

func newStemcell(filePath string) (sc *stemcell, err error) {
	sc = &stemcell{}
	sc.file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}

	sc.gz, err = gzip.NewReader(sc.file)
	if err != nil {
		sc.file.Close()
		return nil, err
	}

	return sc, nil
}

type stemcell struct {
	file *os.File
	gz   *gzip.Reader
	tr   *tar.Reader
}

func (s *stemcell) Close() (err error) {
	err = s.gz.Close()
	err = s.file.Close()
	return
}

func (s *stemcell) Read(p []byte) (n int, err error) {
	return s.gz.Read(p)
}
