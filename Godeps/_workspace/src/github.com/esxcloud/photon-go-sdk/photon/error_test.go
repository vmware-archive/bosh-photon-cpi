package photon

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrorTesting", func() {
	var (
		server *testServer
		client *Client
	)

	BeforeEach(func() {
		server, client = testSetup()
	})

	AfterEach(func() {
		server.Close()
	})

	It("TaskError", func() {
		// Unit test only
		if isIntegrationTest() {
			return
		}
		task := &Task{ID: "fake-id", State: "ERROR", Operation: "fake-op"}
		server.SetResponseJson(200, task)
		task, err := client.Tasks.Wait(task.ID)
		taskErr, ok := err.(TaskError)
		Expect(ok).ShouldNot(BeNil())
		Expect(taskErr.ID).Should(Equal(task.ID))
	})

	It("TaskTimeoutError", func() {
		// Unit test only
		if isIntegrationTest() {
			return
		}
		client.options.TaskPollTimeout = 1 * time.Second
		task := &Task{ID: "fake-id", State: "QUEUED", Operation: "fake-op"}
		server.SetResponseJson(200, task)
		task, err := client.Tasks.Wait(task.ID)
		taskErr, ok := err.(TaskTimeoutError)
		Expect(ok).ShouldNot(BeNil())
		Expect(taskErr.ID).Should(Equal(task.ID))
	})

	It("HttpError", func() {
		// Unit test only
		if isIntegrationTest() {
			return
		}
		client.options.TaskPollTimeout = 1 * time.Second
		task := &Task{ID: "fake-id", State: "QUEUED", Operation: "fake-op"}
		server.SetResponseJson(500, "server error")
		task, err := client.Tasks.Wait(task.ID)
		taskErr, ok := err.(HttpError)
		Expect(ok).ShouldNot(BeNil())
		Expect(taskErr.StatusCode).Should(Equal(500))
	})
})
