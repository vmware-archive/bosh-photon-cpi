package main

import (
	"github.com/esxcloud/bosh-esxcloud-cpi/cmd"
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"runtime"
)

var _ = Describe("AgentEnv", func() {
	type testNetwork struct {
		IP string `json:"ip"`
	}

	var (
		ctx      *cpi.Context
		networks []interface{}
		env      map[string]interface{}
		runner   cmd.Runner
	)

	BeforeEach(func() {
		if runtime.GOOS == "linux" {
			runner = cmd.NewRunner()
		} else {
			runner = &fakeRunner{map[string]string{
				"file": "ISO 9660 CD-ROM filesystem data",
			}}
		}
		ctx = &cpi.Context{
			Config: &cpi.Config{
				Agent: &cpi.AgentConfig{Mbus: "fake-mbus", NTP: []string{"fake-ntp"}},
			},
			Runner: runner,
		}
		env = map[string]interface{}{"prop1": "value1", "prop2": 123}
		nw := []testNetwork{testNetwork{"fake-ip-1"}, testNetwork{"fake-ip-2"}}
		networks = make([]interface{}, len(nw))
		for i, v := range nw {
			networks[i] = interface{}(v)
		}
	})

	// This test requires genisoimage to truly verify ISO creation. On Linux, a real
	// cmd.Runner is used and commands are really executed. On other platforms, it's mocked.
	It("Successfully creates an ISO", func() {
		env := createAgentEnv(ctx, "agent-id", "vm-id", "vm-name", networks, env)
		iso, err := createEnvISO(env, runner)
		defer os.Remove(iso)

		Expect(err).Should(BeNil())

		// Verify we have produced a valid ISO by checking the output of "file <iso>"
		out, err := runner.Run("file", iso)
		outStr := string(out[:])

		Expect(err).Should(BeNil())
		Expect(outStr).Should(ContainSubstring("ISO 9660 CD-ROM"))
	})
})
