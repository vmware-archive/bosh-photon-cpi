package main

import (
	"github.com/esxcloud/bosh-esxcloud-cpi/cpi"
	. "github.com/esxcloud/bosh-esxcloud-cpi/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM metadata", func() {
	It("set_vm_metadata is a nop", func() {
		actions := map[string]cpi.ActionFn{
			"set_vm_metadata": SetVmMetadata,
		}
		res, err := GetResponse(dispatch(nil, actions, "set_vm_metadata", nil))
		Expect(res.Result).To(BeNil())
		Expect(err).To(BeNil())
	})
})
