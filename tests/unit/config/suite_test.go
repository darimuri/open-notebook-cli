package config_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/darimuri/open-notebook-cli/internal/config"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	Describe("Load", func() {
		It("returns default values when no config file exists", func() {
			cfg, err := config.Load("")
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.APIURL).To(Equal("http://localhost:8080"))
			Expect(cfg.Output).To(Equal("table"))
		})
	})
})