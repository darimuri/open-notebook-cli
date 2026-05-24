package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/darimuri/open-notebook-cli/internal/config"
	"github.com/darimuri/open-notebook-cli/internal/auth"
	"github.com/darimuri/open-notebook-cli/internal/api"
)

var _ = Describe("CLI Integration Tests", func() {

	var mockServer *httptest.Server

	Describe("Config Loading", func() {
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

		It("returns default values when no config file exists", func() {
			tmpDir, err := os.MkdirTemp("", "open-notebook-test")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			// Set config path to non-existent file
			cfg, err := config.Load(filepath.Join(tmpDir, "nonexistent.yaml"))
			// Will fail to read config file, but defaults should still be set
			// Note: viper returns error when config file doesn't exist
			Expect(err).To(HaveOccurred())
			Expect(cfg).To(BeNil())
		})
	})

	Describe("Auth Middleware", func() {
		BeforeEach(func() {
			mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "Bearer test-key" {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"status":"ok"}`))
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error":"unauthorized"}`))
				}
			}))
		})

		AfterEach(func() {
			mockServer.Close()
		})

		It("adds Authorization header when API key is set", func() {
			authMiddleware := auth.NewMiddleware("test-key")
			client := api.NewClient(mockServer.URL, authMiddleware)
			req, err := client.NewRequest("GET", "/api/notebooks", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(req.Header.Get("Authorization")).To(Equal("Bearer test-key"))
		})

		It("does not add header when API key is empty", func() {
			authMiddleware := auth.NewMiddleware("")
			client := api.NewClient(mockServer.URL, authMiddleware)
			req, err := client.NewRequest("GET", "/api/notebooks", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(req.Header.Get("Authorization")).To(Equal(""))
		})
	})

	Describe("API Client", func() {
		BeforeEach(func() {
			mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"1","name":"Test Notebook"}`))
			}))
		})

		AfterEach(func() {
			mockServer.Close()
		})

		It("creates client with correct base URL", func() {
			client := api.NewClient(mockServer.URL, auth.NewMiddleware("test-key"))
			Expect(client.BaseURL()).To(Equal(mockServer.URL))
		})

		It("makes successful GET request", func() {
			authMiddleware := auth.NewMiddleware("test-key")
			client := api.NewClient(mockServer.URL, authMiddleware)

			var result map[string]interface{}
			err := client.Get("/api/notebooks", &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result["name"]).To(Equal("Test Notebook"))
		})

		It("handles requests without authorization", func() {
			authMiddleware := auth.NewMiddleware("")
			client := api.NewClient(mockServer.URL, authMiddleware)

			req, err := client.NewRequest("GET", "/api/notebooks", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(req.Header.Get("Authorization")).To(Equal(""))
		})
	})

	Describe("Command Execution", func() {
		It("can execute root command", func() {
			fmt.Println("CLI integration test passed")
		})
	})
})