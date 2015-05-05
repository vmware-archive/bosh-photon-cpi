package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestDispatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dispatch Suite")
	RunSpecs(t, "Stemcell Suite")
}
