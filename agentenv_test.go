package main

import (
	"github.com/esxcloud/bosh-esxcloud-cpi/cmd"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	"github.com/esxcloud/bosh-esxcloud-cpi/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("AgentEnv", func() {
	var (
		ctx      *cpi.Context
		networks map[string]interface{}
		env      map[string]interface{}
		runner   cmd.Runner
	)

	BeforeEach(func() {
		runner = cmd.NewRunner()
		ctx = &cpi.Context{
			Config: &cpi.Config{
				Agent: &cpi.AgentConfig{Mbus: "fake-mbus", NTP: []string{"fake-ntp"}},
			},
			Runner: runner,
			Logger: logger.New(),
		}
		env = map[string]interface{}{"prop1": "value1", "prop2": 123}
		networks = map[string]interface{}{"default": map[string]interface{}{}}
	})

	// This test requires genisoimage to truly verify ISO creation. On Linux, a real
	// cmd.Runner is used and commands are really executed. On other platforms, it's mocked.
	It("Successfully creates an ISO", func() {
		env := createAgentEnv(ctx, "agent-id", "vm-id", "vm-name", networks, env)
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
			env := createAgentEnv(ctx, "agent-id", vmID, "vm-name", map[string]interface{}{}, map[string]interface{}{})
			err := putAgentEnvMetadata(vmID, env)
			Expect(err).ToNot(HaveOccurred())
			env2, err := getAgentEnvMetadata(vmID)
			Expect(err).ToNot(HaveOccurred())
			Expect(env2).Should(Equal(env))
		})
	})
})
