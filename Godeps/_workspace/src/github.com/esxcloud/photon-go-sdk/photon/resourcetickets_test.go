package photon

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResourceTicket", func() {
	var (
		server   *testServer
		client   *Client
		tenantID string
	)

	BeforeEach(func() {
		server, client = testSetup()
		tenantID = createTenant(server, client)
	})

	AfterEach(func() {
		cleanTenants(client)
		server.Close()
	})

	Describe("GetResourceTicketTasks", func() {
		It("GetTasks returns a completed task", func() {
			mockTask := createMockTask("CREATE_RESOURCE_TICKET", "COMPLETED")
			mockTask.Entity.ID = "mock-task-id"
			server.SetResponseJson(200, mockTask)
			spec := &ResourceTicketCreateSpec{
				Name:   randomString(10),
				Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 16, Key: "vm.memory"}},
			}
			task, err := client.Tenants.CreateResourceTicket(tenantID, spec)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_RESOURCE_TICKET"))
			Expect(task.State).Should(Equal("COMPLETED"))

			server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
			taskList, err := client.ResourceTickets.GetTasks(task.Entity.ID, &TaskGetOptions{})
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(taskList).ShouldNot(BeNil())
			Expect(taskList.Items).Should(ContainElement(*task))
		})
	})
})
