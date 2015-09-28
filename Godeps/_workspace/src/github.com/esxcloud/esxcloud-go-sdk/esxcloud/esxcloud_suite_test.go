package esxcloud_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEsxcloud(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Go SDK Suite")
}
