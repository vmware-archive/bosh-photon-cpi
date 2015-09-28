package esxcloud

import (
. "github.com/onsi/ginkgo"
. "github.com/onsi/gomega"
)

var _ = Describe("Status", func() {
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

	// Simple preliminary test. Make sure status API correctly deserializes the response
	It("GetStatus200", func() {
		expectedStruct := Status{"READY", []Component{{"chairman", "", "READY"}, {"housekeeper", "", "READY"}}}
		server.SetResponseJson(200, expectedStruct)

		status, err := client.Status.Get()
		GinkgoT().Log(err)
		Expect(err).Should(BeNil())
		Expect(status.Status).Should(Equal(expectedStruct.Status))
		Expect(status.Components).ShouldNot(HaveLen(1))
	})
})
