package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestDispatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = AfterSuite(func() {
	// Print out temp dir as a helpful reference to find logs
	fmt.Printf("\nLog directory is under: %s\n", os.TempDir())
})
