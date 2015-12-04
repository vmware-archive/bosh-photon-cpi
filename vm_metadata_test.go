package main

import (
	"github.com/esxcloud/bosh-photon-cpi/cpi"
	"github.com/esxcloud/bosh-photon-cpi/logger"
	. "github.com/esxcloud/bosh-photon-cpi/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM metadata", func() {
	var (
		ctx *cpi.Context
	)

	BeforeEach(func() {
		ctx = &cpi.Context{
			Logger: logger.New(),
		}
	})

	It("set_vm_metadata is a nop", func() {
		actions := map[string]cpi.ActionFn{
			"set_vm_metadata": SetVmMetadata,
		}
		res, err := GetResponse(dispatch(ctx, actions, "set_vm_metadata", nil))
		Expect(res.Result).To(BeNil())
		Expect(err).To(BeNil())
		Expect(res.Log).ShouldNot(BeEmpty())
	})
})
