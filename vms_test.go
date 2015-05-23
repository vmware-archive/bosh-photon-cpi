package main

import (
	"github.com/esxcloud/bosh-esxcloud-cpi/cmd"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	. "github.com/esxcloud/bosh-esxcloud-cpi/mocks"
	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"runtime"
)

var _ = Describe("VMs", func() {
	var (
		server *httptest.Server
		ctx    *cpi.Context
		projID string
	)

	BeforeEach(func() {
		server = NewMockServer()
		var runner cmd.Runner
		if runtime.GOOS == "linux" {
			runner = cmd.NewRunner()
		} else {
			runner = &fakeRunner{}
		}

		Activate(true)
		httpClient := &http.Client{Transport: DefaultMockTransport}
		ctx = &cpi.Context{
			Client: ec.NewTestClient(server.URL, httpClient),
			Config: &cpi.Config{
				ESXCloud: &cpi.ESXCloudConfig{
					Target:    server.URL,
					ProjectID: "fake-project-id",
				},
				Agent: &cpi.AgentConfig{Mbus: "fake-mbus", NTP: []string{"fake-ntp"}},
			},
			Runner: runner,
		}

		projID = ctx.Config.ESXCloud.ProjectID
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("CreateVM", func() {
		It("should return ID of created VM", func() {
			createTask := &ec.Task{Operation: "CREATE_VM", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			completedTask := &ec.Task{Operation: "CREATE_VM", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			isoTask := &ec.Task{Operation: "ATTACH_ISO", State: "QUEUED", ID: "fake-iso-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			isoCompletedTask := &ec.Task{Operation: "ATTACH_ISO", State: "COMPLETED", ID: "fake-iso-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			onTask := &ec.Task{Operation: "START_VM", State: "QUEUED", ID: "fake-on-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			onCompletedTask := &ec.Task{Operation: "START_VM", State: "COMPLETED", ID: "fake-on-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			RegisterResponder(
				"POST",
				server.URL+"/v1/projects/"+projID+"/vms",
				CreateResponder(200, ToJson(createTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+createTask.ID,
				CreateResponder(200, ToJson(completedTask)))
			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/"+createTask.Entity.ID+"/attach_iso",
				CreateResponder(200, ToJson(isoTask)))
			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/"+createTask.Entity.ID+"/operations",
				CreateResponder(200, ToJson(onTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+isoTask.ID,
				CreateResponder(200, ToJson(isoCompletedTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+onCompletedTask.ID,
				CreateResponder(200, ToJson(onCompletedTask)))

			actions := map[string]cpi.ActionFn{
				"create_vm": CreateVM,
			}
			args := []interface{}{
				"agent-id",
				"fake-stemcell-id",
				map[string]interface{}{
					"vm_flavor":   "fake-flavor",
					"disk_flavor": "fake-flavor",
				}, // cloud_properties
				map[string]interface{}{}, // networks
				[]string{},               // disk_cids
				map[string]interface{}{}, // environment
			}
			res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

			Expect(res.Result).Should(Equal(completedTask.Entity.ID))
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when server returns error", func() {
			createTask := &ec.Task{Operation: "CREATE_VM", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			completedTask := &ec.Task{Operation: "CREATE_VM", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			RegisterResponder(
				"POST",
				server.URL+"/v1/projects/"+projID+"/vms",
				CreateResponder(500, ToJson(createTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+createTask.ID,
				CreateResponder(200, ToJson(completedTask)))

			actions := map[string]cpi.ActionFn{
				"create_vm": CreateVM,
			}
			args := []interface{}{
				"agent-id",
				"fake-stemcell-id",
				map[string]interface{}{
					"vm_flavor":   "fake-flavor",
					"disk_flavor": "fake-flavor",
				}, // cloud_properties
				map[string]interface{}{}, // networks
				[]string{},               // disk_cids
				map[string]interface{}{}, // environment
			}
			res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when cloud_properties has bad property type", func() {
			actions := map[string]cpi.ActionFn{
				"create_vm": CreateVM,
			}
			args := []interface{}{
				"agent-id",
				"fake-stemcell-id",
				map[string]interface{}{
					"vm_flavor":   123,
					"disk_flavor": "fake-flavor",
				}, // cloud_properties
				map[string]interface{}{}, // networks
				[]string{},               // disk_cids
				map[string]interface{}{}, // environment
			}
			res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when cloud_properties has no properties", func() {
			actions := map[string]cpi.ActionFn{
				"create_vm": CreateVM,
			}
			args := []interface{}{"agent-id", "fake-stemcell-id", map[string]interface{}{}}
			res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when cloud_properties is missing", func() {
			actions := map[string]cpi.ActionFn{
				"create_vm": CreateVM,
			}
			args := []interface{}{"agent-id", "fake-stemcell-id"}
			res, err := GetResponse(dispatch(ctx, actions, "create_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("DeleteVM", func() {
		It("should return nothing when successful", func() {
			deleteTask := &ec.Task{Operation: "DELETE_VM", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			completedTask := &ec.Task{Operation: "DELETE_VM", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			offTask := &ec.Task{Operation: "STOP_VM", State: "QUEUED", ID: "fake-off-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			offCompletedTask := &ec.Task{Operation: "STOP_VM", State: "COMPLETED", ID: "fake-off-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			RegisterResponder(
				"DELETE",
				server.URL+"/v1/vms/"+deleteTask.Entity.ID+"?force=true",
				CreateResponder(200, ToJson(deleteTask)))
			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/"+deleteTask.Entity.ID+"/operations",
				CreateResponder(200, ToJson(offTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+deleteTask.ID,
				CreateResponder(200, ToJson(completedTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+offCompletedTask.ID,
				CreateResponder(200, ToJson(offCompletedTask)))

			actions := map[string]cpi.ActionFn{
				"delete_vm": DeleteVM,
			}
			args := []interface{}{"fake-vm-id"}
			res, err := GetResponse(dispatch(ctx, actions, "delete_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when VM not found", func() {
			deleteTask := &ec.Task{Operation: "DELETE_VM", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			completedTask := &ec.Task{Operation: "DELETE_VM", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			RegisterResponder(
				"DELETE",
				server.URL+"/v1/vms/"+deleteTask.Entity.ID+"?force=true",
				CreateResponder(404, ToJson(deleteTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+deleteTask.ID,
				CreateResponder(200, ToJson(completedTask)))

			actions := map[string]cpi.ActionFn{
				"delete_vm": DeleteVM,
			}
			args := []interface{}{"fake-vm-id"}
			res, err := GetResponse(dispatch(ctx, actions, "delete_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given no arguments", func() {
			actions := map[string]cpi.ActionFn{
				"delete_vm": DeleteVM,
			}
			args := []interface{}{}
			res, err := GetResponse(dispatch(ctx, actions, "delete_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given an invalid argument", func() {
			actions := map[string]cpi.ActionFn{
				"delete_vm": DeleteVM,
			}
			args := []interface{}{5}
			res, err := GetResponse(dispatch(ctx, actions, "delete_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Describe("HasVM", func() {
		It("should return true when VM is found", func() {
			vm := &ec.VM{ID: "fake-vm-id"}
			RegisterResponder(
				"GET",
				server.URL+"/v1/vms/"+vm.ID,
				CreateResponder(200, ToJson(vm)))

			actions := map[string]cpi.ActionFn{
				"has_vm": HasVM,
			}
			args := []interface{}{"fake-vm-id"}
			res, err := GetResponse(dispatch(ctx, actions, "has_vm", args))

			Expect(res.Result).Should(Equal(true))
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return false when VM not found", func() {
			vm := &ec.VM{ID: "fake-vm-id"}
			RegisterResponder(
				"GET",
				server.URL+"/v1/vms/"+vm.ID,
				CreateResponder(404, ToJson(vm)))

			actions := map[string]cpi.ActionFn{
				"has_vm": HasVM,
			}
			args := []interface{}{"fake-vm-id"}
			res, err := GetResponse(dispatch(ctx, actions, "has_vm", args))

			Expect(res.Result).Should(Equal(false))
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when server returns error", func() {
			vm := &ec.VM{ID: "fake-vm-id"}
			RegisterResponder(
				"GET",
				server.URL+"/v1/vms/"+vm.ID,
				CreateResponder(500, ToJson(vm)))

			actions := map[string]cpi.ActionFn{
				"has_vm": HasVM,
			}
			args := []interface{}{"fake-vm-id"}
			res, err := GetResponse(dispatch(ctx, actions, "has_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given no arguments", func() {
			actions := map[string]cpi.ActionFn{
				"has_vm": HasVM,
			}
			args := []interface{}{}
			res, err := GetResponse(dispatch(ctx, actions, "has_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given an invalid argument", func() {
			actions := map[string]cpi.ActionFn{
				"has_vm": HasVM,
			}
			args := []interface{}{5}
			res, err := GetResponse(dispatch(ctx, actions, "has_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("RestartVM", func() {
		It("should return nothing when successful", func() {
			restartTask := &ec.Task{Operation: "restart_vm", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			completedTask := &ec.Task{Operation: "restart_vm", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/"+restartTask.Entity.ID+"/operations",
				CreateResponder(200, ToJson(restartTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+restartTask.ID,
				CreateResponder(200, ToJson(completedTask)))

			actions := map[string]cpi.ActionFn{
				"restart_vm": RestartVM,
			}
			args := []interface{}{"fake-vm-id"}
			res, err := GetResponse(dispatch(ctx, actions, "restart_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).Should(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when VM not found", func() {
			restartTask := &ec.Task{Operation: "restart_vm", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}
			completedTask := &ec.Task{Operation: "restart_vm", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-vm-id"}}

			RegisterResponder(
				"POST",
				server.URL+"/v1/vms/"+restartTask.Entity.ID+"/operations",
				CreateResponder(404, ToJson(restartTask)))
			RegisterResponder(
				"GET",
				server.URL+"/v1/tasks/"+restartTask.ID,
				CreateResponder(200, ToJson(completedTask)))

			actions := map[string]cpi.ActionFn{
				"restart_vm": RestartVM,
			}
			args := []interface{}{"fake-vm-id"}
			res, err := GetResponse(dispatch(ctx, actions, "restart_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given no arguments", func() {
			actions := map[string]cpi.ActionFn{
				"restart_vm": RestartVM,
			}
			args := []interface{}{}
			res, err := GetResponse(dispatch(ctx, actions, "restart_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should return an error when given an invalid argument", func() {
			actions := map[string]cpi.ActionFn{
				"restart_vm": RestartVM,
			}
			args := []interface{}{5}
			res, err := GetResponse(dispatch(ctx, actions, "restart_vm", args))

			Expect(res.Result).Should(BeNil())
			Expect(res.Error).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
