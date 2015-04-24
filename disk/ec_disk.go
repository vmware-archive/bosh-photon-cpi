package disk

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
)

const esxcloudDiskLogTag = "ESXCloudDisk"

type ECDisk struct {
	id   			string
	esxCloudClient	ec.Client
	logger 			boshlog.Logger
}

func NewECDisk(id string, client ec.Client, logger boshlog.Logger) *ECDisk {
	return &ECDisk{id: id, esxCloudClient: client, logger: logger}
}

func (d *ECDisk) ID() string { return d.id }

func (d *ECDisk) Delete() error {
	d.logger.Debug(esxcloudDiskLogTag, "Deleting disk '%s'", d.id)

	deleteTask, err := d.esxCloudClient.Disks.Delete(d.id, false)
	if err != nil {
		return bosherr.WrapErrorf(err, "Fail to create a task to delete disk '%s'", d.id)
	}
	if deleteTask == nil {
		return bosherr.WrapErrorf(bosherr.Error("Fail to delete disk '%s'."),
									"No task received from API when a task was expected, " +
		                    		"but no error was received, which should not happen",
									d.id)
	}

	waitTask, err := d.esxCloudClient.Tasks.Wait(deleteTask.ID)
	if err != nil {
		return bosherr.WrapErrorf(err, "Task failed: delete disk '%s'", d.id)
	}
	if waitTask == nil {
		return bosherr.WrapErrorf(bosherr.Error("Fail to delete disk '%s' when wait for it."),
									"No task received from API when a task was expected, " +
									"but no error was received, which should not happen.",
									d.id)
	}
	if waitTask.State != "COMPLETED" {
		return bosherr.WrapErrorf(bosherr.Error("Delete disk task is not complete"),
									"Task of delete disk '%s' is not Completed", d.id)
	}

	return nil
}
