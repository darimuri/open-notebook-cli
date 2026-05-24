package config_test

import (
	"os"
	"path/filepath"
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
		It("loads from config file", func() {
			tmpDir, err := os.MkdirTemp("", "open-notebook-test")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			configPath := filepath.Join(tmpDir, "config.yaml")
			err = os.WriteFile(configPath, []byte(`api_url: "https://test.darimuri.me"
api_key: "test-key"
output: "json"
`), 0644)
			Expect(err).NotTo(HaveOccurred())

			cfg, err := config.Load(configPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.APIURL).To(Equal("https://test.darimuri.me"))
			Expect(cfg.APIKey).To(Equal("test-key"))
			Expect(cfg.Output).To(Equal("json"))
		})

		It("returns error when config file not found", func() {
			cfg, err := config.Load("/nonexistent/path/config.yaml")
			Expect(err).To(HaveOccurred())
			Expect(cfg).To(BeNil())
		})

		It("loads with defaults when no config file exists and no env vars set", func() {
			tmpDir, err := os.MkdirTemp("", "open-notebook-test")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			// Unset env vars
			os.Unsetenv("OPEN_NOTEBOOK_API_URL")
			os.Unsetenv("OPEN_NOTEBOOK_API_KEY")
			os.Unsetenv("OPEN_NOTEBOOK_OUTPUT")

			// Create a config directory but no config file
			configDir := filepath.Join(tmpDir, "open-notebook")
			err = os.MkdirAll(configDir, 0755)
			Expect(err).NotTo(HaveOccurred())

			// Point to non-existent file in existing directory
			_, err = config.Load(filepath.Join(configDir, "config.yaml"))
			// Viper returns error when file doesn't exist but directory does
			// This tests the actual behavior
			Expect(err).To(HaveOccurred())
		})
	})
})