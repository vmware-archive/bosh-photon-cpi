package disk_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	. "github.com/esxcloud/bosh-esxcloud-cpi/disk"
	. "github.com/esxcloud/bosh-esxcloud-cpi/mocks"
	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
)

var _ = Describe("ECDisk", func() {
	var (
		fc   *ec.Client
		disk *ECDisk
		uri  string

		server *httptest.Server
	)

	BeforeEach(func() {
		server = NewMockServer()
		uri = server.URL

		Activate(true)
		httpClient := &http.Client{Transport: DefaultMockTransport}
		fc = ec.NewTestClient(uri, httpClient)

		logger := boshlog.NewLogger(boshlog.LevelNone)
		disk = NewECDisk("fake-disk-id", *fc, logger)
	})

	AfterEach(func() {
		server.Close()
	})

	It("deletes a persistant disk successfully", func() {
		delete_task := NewMockTask("DELETE_DISK", "QUEUED", "fake-delete-task-id")
		pull_task := NewMockTask("DELETE_DISK", "COMPLETED", "")

		RegisterResponder(
			"DELETE",
			uri+"/v1/disks/fake-disk-id?force=false",
			CreateResponder(200, ToJson(delete_task)))

		RegisterResponder(
			"GET",
			uri+"/v1/tasks/fake-delete-task-id",
			CreateResponder(200, ToJson(pull_task)))

		err := disk.Delete()
		Expect(err).ToNot(HaveOccurred())
	})
})
