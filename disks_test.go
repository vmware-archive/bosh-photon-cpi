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

var _ = Describe("Disk", func() {
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
					Target:     server.URL,
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

	Describe("CreateDisk", func() {
		It("returns a disk ID", func() {
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
		It("returns an error when size is too small", func() {
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
		It("returns an error when apife returns a 500", func() {
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
		It("should return an error when given no arguments", func() {
			actions := map[string]cpi.ActionFn{
				"create_disk": CreateDisk,
			}
			args := []interface{}{}
			res, err := GetResponse(dispatch(ctx, actions, "create_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given an invalid argument", func() {
			actions := map[string]cpi.ActionFn{
				"create_disk": CreateDisk,
			}
			args := []interface{}{"not-an-int"}
			res, err := GetResponse(dispatch(ctx, actions, "create_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("DeleteDisk", func() {
		It("returns nothing", func() {
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
		It("returns an error when apife returns 404", func() {
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
		It("should return an error when given no arguments", func() {
			actions := map[string]cpi.ActionFn{
				"delete_disk": DeleteDisk,
			}
			args := []interface{}{}
			res, err := GetResponse(dispatch(ctx, actions, "delete_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given an invalid argument", func() {
			actions := map[string]cpi.ActionFn{
				"delete_disk": DeleteDisk,
			}
			args := []interface{}{5}
			res, err := GetResponse(dispatch(ctx, actions, "delete_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("HasDisk", func() {
		It("returns true for HasDisk when disk exists", func() {
			disk := &ec.PersistentDisk{Flavor: "persistent-disk", ID: "fake-disk-id"}

			RegisterResponder(
				"GET",
				server.URL+"/v1/disks/"+disk.ID,
				CreateResponder(200, ToJson(disk)))

			actions := map[string]cpi.ActionFn{
				"has_disk": HasDisk,
			}
			args := []interface{}{"fake-disk-id"}
			res, err := GetResponse(dispatch(ctx, actions, "has_disk", args))

			Expect(res.Result).Should(Equal(true))
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("returns false for HasDisk when disk does not exists", func() {
			disk := &ec.PersistentDisk{Flavor: "persistent-disk", ID: "fake-disk-id"}

			RegisterResponder(
				"GET",
				server.URL+"/v1/disks/"+disk.ID,
				CreateResponder(404, ToJson(disk)))

			actions := map[string]cpi.ActionFn{
				"has_disk": HasDisk,
			}
			args := []interface{}{"fake-disk-id"}
			res, err := GetResponse(dispatch(ctx, actions, "has_disk", args))

			Expect(res.Result).Should(Equal(false))
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("returns an error for HasDisk when server returns error", func() {
			disk := &ec.PersistentDisk{Flavor: "persistent-disk", ID: "fake-disk-id"}

			RegisterResponder(
				"GET",
				server.URL+"/v1/disks/"+disk.ID,
				CreateResponder(500, ToJson(disk)))

			actions := map[string]cpi.ActionFn{
				"has_disk": HasDisk,
			}
			args := []interface{}{"fake-disk-id"}
			res, err := GetResponse(dispatch(ctx, actions, "has_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given no arguments", func() {
			actions := map[string]cpi.ActionFn{
				"has_disk": HasDisk,
			}
			args := []interface{}{}
			res, err := GetResponse(dispatch(ctx, actions, "has_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given an invalid argument", func() {
			actions := map[string]cpi.ActionFn{
				"has_disk": HasDisk,
			}
			args := []interface{}{5}
			res, err := GetResponse(dispatch(ctx, actions, "has_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("GetDisks", func() {
		It("returns a list of VM IDs that a disk is attached to", func() {
			list := &ec.DiskList{
				[]ec.PersistentDisk{
					ec.PersistentDisk{ID: "disk-1", VMs: []string{"vm-1", "vm-2"}},
					ec.PersistentDisk{ID: "disk-2", VMs: []string{"vm-2", "vm-3"}},
					ec.PersistentDisk{ID: "disk-3", VMs: []string{"vm-4", "vm-5"}},
					ec.PersistentDisk{ID: "disk-4", VMs: []string{"vm-2", "vm-4"}},
				},
			}
			// Disks on vm-2
			matchedList := []interface{}{"disk-4", "disk-2", "disk-1"}

			RegisterResponder(
				"GET",
				server.URL+"/v1/projects/"+projID+"/disks",
				CreateResponder(200, ToJson(list)))

			actions := map[string]cpi.ActionFn{
				"get_disks": GetDisks,
			}
			args := []interface{}{"vm-2"}
			res, err := GetResponse(dispatch(ctx, actions, "get_disks", args))

			Expect(res.Result).Should(ConsistOf(matchedList))
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("returns an empty list if no disks are attached to VM", func() {
			list := &ec.DiskList{
				[]ec.PersistentDisk{
					ec.PersistentDisk{ID: "disk-1", VMs: []string{"other-vm"}},
				},
			}
			matchedList := []interface{}{}

			RegisterResponder(
				"GET",
				server.URL+"/v1/projects/"+projID+"/disks",
				CreateResponder(200, ToJson(list)))

			actions := map[string]cpi.ActionFn{
				"get_disks": GetDisks,
			}
			args := []interface{}{"vm-2"}
			res, err := GetResponse(dispatch(ctx, actions, "get_disks", args))

			Expect(res.Result).Should(ConsistOf(matchedList))
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("returns an error when server returns error", func() {
			list := &ec.DiskList{[]ec.PersistentDisk{}}

			RegisterResponder(
				"GET",
				server.URL+"/v1/projects/"+projID+"/disks",
				CreateResponder(500, ToJson(list)))

			actions := map[string]cpi.ActionFn{
				"get_disks": GetDisks,
			}
			args := []interface{}{"vm-2"}
			res, err := GetResponse(dispatch(ctx, actions, "get_disks", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Describe("AttachDisk", func() {
		It("returns nothing when attach succeeds", func() {
			attachTask := &ec.Task{Operation: "ATTACH_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
			completedTask := &ec.Task{Operation: "ATTACH_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/fake-vm-id/attach_disk",
				CreateResponder(200, ToJson(attachTask)))

			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+attachTask.ID,
				CreateResponder(200, ToJson(completedTask)))

			actions := map[string]cpi.ActionFn{
				"attach_disk": AttachDisk,
			}
			args := []interface{}{"fake-vm-id", "fake-disk-id"}
			res, err := GetResponse(dispatch(ctx, actions, "attach_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("returns an error when VM not found", func() {
			attachTask := &ec.Task{Operation: "ATTACH_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
			completedTask := &ec.Task{Operation: "ATTACH_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/fake-vm-id/attach_disk",
				CreateResponder(404, ToJson(attachTask)))

			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+attachTask.ID,
				CreateResponder(200, ToJson(completedTask)))

			actions := map[string]cpi.ActionFn{
				"attach_disk": AttachDisk,
			}
			args := []interface{}{"fake-vm-id", "fake-disk-id"}
			res, err := GetResponse(dispatch(ctx, actions, "attach_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Describe("DetachDisk", func() {
		It("returns nothing when detach succeeds", func() {
			attachTask := &ec.Task{Operation: "DETACH_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
			completedTask := &ec.Task{Operation: "DETACH_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/fake-vm-id/detach_disk",
				CreateResponder(200, ToJson(attachTask)))

			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+attachTask.ID,
				CreateResponder(200, ToJson(completedTask)))

			actions := map[string]cpi.ActionFn{
				"detach_disk": DetachDisk,
			}
			args := []interface{}{"fake-vm-id", "fake-disk-id"}
			res, err := GetResponse(dispatch(ctx, actions, "detach_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("returns an error when VM not found", func() {
			attachTask := &ec.Task{Operation: "DETACH_DISK", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}
			completedTask := &ec.Task{Operation: "DETACH_DISK", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-disk-id"}}

			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/fake-vm-id/detach_disk",
				CreateResponder(404, ToJson(attachTask)))

			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+attachTask.ID,
				CreateResponder(200, ToJson(completedTask)))

			actions := map[string]cpi.ActionFn{
				"attach_disk": DetachDisk,
			}
			args := []interface{}{"fake-vm-id", "fake-disk-id"}
			res, err := GetResponse(dispatch(ctx, actions, "attach_disk", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
