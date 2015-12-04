package photon_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPhoton(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Go SDK Suite")
}
