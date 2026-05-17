package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/darimuri/open-notebook-cli/internal/api"
	"github.com/darimuri/open-notebook-cli/internal/auth"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("Client", func() {
	Describe("NewClient", func() {
		It("creates client with correct base URL", func() {
			client := api.NewClient("http://localhost:8080", auth.NewMiddleware("test-key"))
			Expect(client.BaseURL()).To(Equal("http://localhost:8080"))
		})
	})

	Describe("NewRequest", func() {
		It("sets auth header on request", func() {
			client := api.NewClient("http://localhost:8080", auth.NewMiddleware("test-key"))
			req, err := client.NewRequest("GET", "/api/notebooks", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(req.Header.Get("Authorization")).To(Equal("Bearer test-key"))
		})

		It("sets Content-Type header", func() {
			client := api.NewClient("http://localhost:8080", auth.NewMiddleware(""))
			req, err := client.NewRequest("POST", "/api/notebooks", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
		})
	})
})