package main

import (
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	. "github.com/esxcloud/bosh-esxcloud-cpi/mocks"
	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Stemcell", func() {
	var (
		server *httptest.Server
		ctx    *cpi.Context
		projID string
	)

	BeforeEach(func() {
		server = NewMockServer()

		Activate(true)
		httpClient := &http.Client{Transport: DefaultMockTransport}
		ctx = &cpi.Context{
			Client: ec.NewTestClient(server.URL, httpClient),
			Config: &cpi.Config{
				ESXCloud: &cpi.ESXCloudConfig{
					APIFE:      server.URL,
					DiskFlavor: "test-disk-flavor",
					ProjectID:  "fake-project-id",
				},
			},
		}

		projID = ctx.Config.ESXCloud.ProjectID
	})

	AfterEach(func() {
		server.Close()
	})

	It("returns a disk ID for Create", func() {
		createTask := &ec.Task{Operation: "CREATE_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
		completedTask := &ec.Task{Operation: "CREATE_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

		RegisterResponder(
			"POST",
			server.URL+"/v1/projects/"+projID+"/disks",
			CreateResponder(200, ToJson(createTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+createTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]cpi.ActionFn{
			"create_disk": CreateDisk,
		}
		args := []interface{}{2500, "fake-vm-id"}
		res, err := GetResponse(dispatch(ctx, actions, "create_disk", args))

		Expect(res.Result).Should(Equal(completedTask.Entity.ID))
		Expect(res.Error).Should(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns an error when Create size is too small", func() {
		createTask := &ec.Task{Operation: "CREATE_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
		completedTask := &ec.Task{Operation: "CREATE_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

		RegisterResponder(
			"POST",
			server.URL+"/v1/projects/"+projID+"/disks",
			CreateResponder(200, ToJson(createTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+createTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]cpi.ActionFn{
			"create_disk": CreateDisk,
		}
		args := []interface{}{0, "fake-vm-id"}
		res, err := GetResponse(dispatch(ctx, actions, "create_disk", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns an error when Create size is too small", func() {
		createTask := &ec.Task{Operation: "CREATE_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
		completedTask := &ec.Task{Operation: "CREATE_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

		RegisterResponder(
			"POST",
			server.URL+"/v1/projects/"+projID+"/disks",
			CreateResponder(200, ToJson(createTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+createTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]cpi.ActionFn{
			"create_disk": CreateDisk,
		}
		args := []interface{}{0, "fake-vm-id"}
		res, err := GetResponse(dispatch(ctx, actions, "create_disk", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns an error for Create when apife returns a 500", func() {
		createTask := &ec.Task{Operation: "CREATE_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
		completedTask := &ec.Task{Operation: "CREATE_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

		RegisterResponder(
			"POST",
			server.URL+"/v1/projects/"+projID+"/disks",
			CreateResponder(500, ToJson(createTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+createTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]cpi.ActionFn{
			"create_disk": CreateDisk,
		}
		args := []interface{}{2500, "fake-vm-id"}
		res, err := GetResponse(dispatch(ctx, actions, "create_disk", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns nothing for Delete", func() {
		deleteTask := &ec.Task{Operation: "DELETE_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
		completedTask := &ec.Task{Operation: "DELETE_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

		RegisterResponder(
			"DELETE",
			server.URL+"/v1/disks/"+deleteTask.Entity.ID+"?force=true",
			CreateResponder(200, ToJson(deleteTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+deleteTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]cpi.ActionFn{
			"delete_disk": DeleteDisk,
		}
		args := []interface{}{"fake-disk-id"}
		res, err := GetResponse(dispatch(ctx, actions, "delete_disk", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).Should(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns an error for Delete when apife returns 404", func() {
		deleteTask := &ec.Task{Operation: "DELETE_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
		completedTask := &ec.Task{Operation: "DELETE_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

		RegisterResponder(
			"DELETE",
			server.URL+"/v1/disks/"+deleteTask.Entity.ID+"?force=true",
			CreateResponder(404, ToJson(deleteTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+deleteTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]cpi.ActionFn{
			"delete_disk": DeleteDisk,
		}
		args := []interface{}{"fake-disk-id"}
		res, err := GetResponse(dispatch(ctx, actions, "delete_disk", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
})
