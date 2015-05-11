package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	"io"
	"os"
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
	task, err := ctx.Client.Images.Create(stemcell)
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

	sc.tr = tar.NewReader(sc.gz)

	found := false
	for {
		header, err := sc.tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if header.Typeflag == tar.TypeReg {
			if header.Name == "image" {
				found = true
				break
			}
		}
	}
	if !found {
		err = errors.New(fmt.Sprintf("Could not find entry for OVA image in stemcell at path '%s'", filePath))
		return nil, err
	}
	sc.innerGz, err = gzip.NewReader(sc.tr)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

type stemcell struct {
	file    *os.File
	gz      *gzip.Reader
	tr      *tar.Reader
	innerGz *gzip.Reader
}

func (s *stemcell) Close() (err error) {
	err = s.innerGz.Close()
	err = s.gz.Close()
	err = s.file.Close()
	return
}

func (s *stemcell) Read(p []byte) (n int, err error) {
	return s.innerGz.Read(p)
}
