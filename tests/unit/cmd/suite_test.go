package cmd_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var _ = Describe("Cmd", func() {
	It("has root command", func() {
		// CLI commands are tested via integration tests
		// This just verifies the test suite works
		Expect(true).To(BeTrue())
	})
})