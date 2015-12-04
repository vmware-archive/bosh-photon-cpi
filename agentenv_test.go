package main

import (
	"github.com/esxcloud/bosh-photon-cpi/cmd"
	"github.com/esxcloud/bosh-photon-cpi/cpi"
	"github.com/esxcloud/bosh-photon-cpi/logger"
	. "github.com/esxcloud/bosh-photon-cpi/mocks"
	ec "github.com/esxcloud/photon-go-sdk/photon"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
)

var _ = Describe("AgentEnv", func() {
	var (
		ctx    *cpi.Context
		env    *cpi.AgentEnv
		runner cmd.Runner
		server *httptest.Server
	)

	BeforeEach(func() {
		server = NewMockServer()
		runner = cmd.NewRunner()
		httpClient := &http.Client{Transport: DefaultMockTransport}
		ctx = &cpi.Context{
			Client: ec.NewTestClient(server.URL, nil, httpClient),
			Config: &cpi.Config{
				Photon: &cpi.PhotonConfig{
					Target:    server.URL,
					ProjectID: "fake-project-id",
				},
				Agent: &cpi.AgentConfig{Mbus: "fake-mbus", NTP: []string{"fake-ntp"}},
			},
			Runner: runner,
			Logger: logger.New(),
		}
		env = &cpi.AgentEnv{AgentID: "agent-id", VM: cpi.VMSpec{Name: "vm-name", ID: "vm-id"}}
	})

	It("Successfully creates an ISO", func() {
		iso, err := createEnvISO(env, runner)
		defer os.Remove(iso)

		Expect(err).Should(BeNil(), "Test requires mkisofs, install with 'brew install cdrtools' on Mac")

		// Verify we have produced a valid ISO by checking the output of "file <iso>"
		out, err := runner.Run("file", iso)
		outStr := string(out[:])

		Expect(err).Should(BeNil())
		Expect(outStr).Should(ContainSubstring("ISO 9660 CD-ROM"))
	})

	Describe("Metadata", func() {
		It("successfully puts and gets agent env data", func() {
			vmID := "fake-vm-id"
			metadataTask := &ec.Task{State: "COMPLETED"}
			vm := &ec.VM{
				ID:       vmID,
				Metadata: map[string]string{"bosh-cpi": GetEnvMetadata(env)},
			}

			RegisterResponder(
				"POST",
				server.URL+"/vms/"+vmID+"/set_metadata",
				CreateResponder(200, ToJson(metadataTask)))
			RegisterResponder(
				"GET",
				server.URL+"/vms/"+vmID,
				CreateResponder(200, ToJson(vm)))

			err := putAgentEnvMetadata(ctx, vmID, env)
			Expect(err).ToNot(HaveOccurred())

			env2, err := getAgentEnvMetadata(ctx, vmID)
			Expect(err).ToNot(HaveOccurred())
			Expect(env2).Should(Equal(env))
		})
	})
})
