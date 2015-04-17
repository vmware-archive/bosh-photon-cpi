package disk

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

)

const esxcloudDiskLogTag = "ESXCloudDisk"

type ECDisk struct {
	id   string
	path string

	fs     boshsys.FileSystem
	logger boshlog.Logger
}

func NewECDisk(
	id string,
	path string,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) ECDisk {
	return ECDisk{id: id, path: path, fs: fs, logger: logger}
}

func (s *ECDisk) ID() string { return s.id }

func (s *ECDisk) Path() string { return s.path }

func (s *ECDisk) Delete() error {
	s.logger.Debug(esxcloudDiskLogTag, "Deleting disk '%s'", s.id)

	err := s.fs.RemoveAll(s.path)
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting disk '%s'", s.path)
	}

	return nil
}
