package inttests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestDispatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ESXCloud Bosh CPI Integration Tests")
}
