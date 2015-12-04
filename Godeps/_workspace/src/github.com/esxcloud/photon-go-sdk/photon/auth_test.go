package photon

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
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

	Describe("GetAuth", func() {
		It("returns auth info", func() {
			expected := &AuthInfo{
				Enabled: false,
				Port:    0,
			}
			server.SetResponseJson(200, expected)
			info, err := client.Auth.Get()
			fmt.Fprintf(GinkgoWriter, "Got auth info: %+v\n", info)
			Expect(info).ShouldNot(BeNil())
			Expect(err).Should(BeNil())
		})
	})
})
