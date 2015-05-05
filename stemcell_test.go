package main

import (
	. "github.com/esxcloud/bosh-esxcloud-cpi/mocks"
	"github.com/esxcloud/bosh-esxcloud-cpi/types"
	ec "github.com/esxcloud/esxcloud-go-sdk/esxcloud"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Stemcell", func() {
	var (
		server *httptest.Server
		ctx    *types.Context
	)

	BeforeEach(func() {
		server = NewMockServer()

		Activate(true)
		httpClient := &http.Client{Transport: DefaultMockTransport}
		ctx = &types.Context{ECClient: ec.NewTestClient(server.URL, httpClient)}
	})

	AfterEach(func() {
		server.Close()
	})

	It("returns a stemcell ID for Create", func() {
		createTask := &ec.Task{Operation: "CREATE_IMAGE", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-image-id"}}
		completedTask := &ec.Task{Operation: "CREATE_IMAGE", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-image-id"}}

		RegisterResponder(
			"POST",
			server.URL+"/v1/images",
			CreateResponder(200, ToJson(createTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+createTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]types.ActionFn{
			"create_stemcell": CreateStemcell,
		}
		args := []interface{}{"./testdata/tty_tiny.ova"}
		res, err := GetResponse(dispatch(ctx, actions, "create_stemcell", args))

		Expect(res.Result).Should(Equal(completedTask.Entity.ID))
		Expect(res.Error).Should(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("returns an error when APIfe returns a 500", func() {
		createTask := &ec.Task{Operation: "CREATE_IMAGE", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-image-id"}}
		completedTask := &ec.Task{Operation: "CREATE_IMAGE", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-image-id"}}

		RegisterResponder(
			"POST",
			server.URL+"/v1/images",
			CreateResponder(500, ToJson(createTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+createTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]types.ActionFn{
			"create_stemcell": CreateStemcell,
		}
		args := []interface{}{"./testdata/tty_tiny.ova"}
		res, err := GetResponse(dispatch(ctx, actions, "create_stemcell", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("returns an error when stemcell file does not exist", func() {
		actions := map[string]types.ActionFn{
			"create_stemcell": CreateStemcell,
		}
		args := []interface{}{"a-file-that-does-not-exist"}
		res, err := GetResponse(dispatch(ctx, actions, "create_stemcell", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("returns nothing for stemcell delete", func() {
		deleteTask := &ec.Task{Operation: "DELETE_IMAGE", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-image-id"}}
		completedTask := &ec.Task{Operation: "DELETE_IMAGE", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-image-id"}}

		RegisterResponder(
			"DELETE",
			server.URL+"/v1/images/"+deleteTask.Entity.ID,
			CreateResponder(200, ToJson(deleteTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+deleteTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]types.ActionFn{
			"delete_stemcell": DeleteStemcell,
		}
		args := []interface{}{deleteTask.Entity.ID}
		res, err := GetResponse(dispatch(ctx, actions, "delete_stemcell", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).Should(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("returns an error for missing stemcell delete", func() {
		deleteTask := &ec.Task{Operation: "DELETE_IMAGE", State: "QUEUED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-image-id"}}
		completedTask := &ec.Task{Operation: "DELETE_IMAGE", State: "COMPLETED", ID: "fake-task-id", Entity: ec.Entity{ID: "fake-image-id"}}

		RegisterResponder(
			"DELETE",
			server.URL+"/v1/images/"+deleteTask.Entity.ID,
			CreateResponder(404, ToJson(deleteTask)))
		RegisterResponder(
			"GET",
			server.URL+"/v1/tasks/"+deleteTask.ID,
			CreateResponder(200, ToJson(completedTask)))

		actions := map[string]types.ActionFn{
			"delete_stemcell": DeleteStemcell,
		}
		args := []interface{}{deleteTask.Entity.ID}
		res, err := GetResponse(dispatch(ctx, actions, "delete_stemcell", args))

		Expect(res.Result).Should(BeNil())
		Expect(res.Error).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
})
