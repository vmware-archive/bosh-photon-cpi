package esxcloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type options struct {
	A int    `urlParam:"a"`
	B string `urlParam:"b"`
}

var _ = Describe("Utils", func() {
	It("GetQueryString", func() {
		opts := &options{5, "a test"}
		query := getQueryString(opts)
		Expect(query).Should(Equal("?a=5&b=a+test"))
	})
})
