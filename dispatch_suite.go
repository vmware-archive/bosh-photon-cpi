package main

import (
	"encoding/json"
	"errors"
	. "github.com/esxcloud/bosh-esxcloud-cpi/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dispatch", func() {
	It("returns a valid bosh JSON response given valid arguments", func() {
		actions := map[string]actionFn{
			"create_vm": createVM,
		}
		args := []interface{}{"fake-agent-id"}
		res, err := getResponse(dispatch(actions, "create_vm", args))

		Expect(res.Result).Should(Equal("fake-vm-id"))
		Expect(res.Error).Should(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns a valid bosh JSON error when given an invalid argument", func() {
		actions := map[string]actionFn{
			"create_vm": createVM,
		}
		args := []interface{}{5}
		res, err := getResponse(dispatch(actions, "create_vm", args))

		Expect(res.Error).ShouldNot(BeNil())
		Expect(res.Error.Type).Should(Equal(CloudError))
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns a valid bosh JSON error when function errors", func() {
		actions := map[string]actionFn{
			"create_vm": createVmError,
		}
		args := []interface{}{"fake-agent-id"}
		res, err := getResponse(dispatch(actions, "create_vm", args))

		Expect(res.Error).ShouldNot(BeNil())
		Expect(res.Error.Type).Should(Equal(CloudError))
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns a valid bosh JSON error when function panics", func() {
		actions := map[string]actionFn{
			"create_vm": createVmPanic,
		}
		args := []interface{}{"fake-agent-id"}
		res, err := getResponse(dispatch(actions, "create_vm", args))

		Expect(res.Error).ShouldNot(BeNil())
		Expect(res.Error.Type).Should(Equal(CloudError))
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("returns a valid bosh JSON error when method not implemented", func() {
		actions := map[string]actionFn{}
		args := []interface{}{"fake-agent-id"}
		res, err := getResponse(dispatch(actions, "create_vm", args))

		Expect(res.Error).ShouldNot(BeNil())
		Expect(res.Error.Type).Should(Equal(NotImplementedError))
		Expect(err).ShouldNot(HaveOccurred())
	})
})

func createVM(args []interface{}) (result interface{}, err error) {
	_, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where agent_id should be")
	}
	return "fake-vm-id", nil
}

func createVmError(args []interface{}) (result interface{}, err error) {
	return nil, errors.New("error occurred")
}

func createVmPanic(args []interface{}) (result interface{}, err error) {
	panic("oh no!")
}

func getResponse(data []byte) (res Response, err error) {
	err = json.Unmarshal(data, &res)
	return
}