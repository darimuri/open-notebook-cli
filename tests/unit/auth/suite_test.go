package auth_test

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/darimuri/open-notebook-cli/internal/auth"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

var _ = Describe("Auth", func() {
	Describe("Middleware", func() {
		It("adds Authorization header when API key is set", func() {
			middleware := auth.NewMiddleware("test-key")
			req, err := makeRequest("GET", "/test")
			Expect(err).NotTo(HaveOccurred())

			middleware.AddAuth(req)

			Expect(req.Header.Get("Authorization")).To(Equal("Bearer test-key"))
		})

		It("does not add header when API key is empty", func() {
			middleware := auth.NewMiddleware("")
			req, err := makeRequest("GET", "/test")
			Expect(err).NotTo(HaveOccurred())

			middleware.AddAuth(req)

			Expect(req.Header.Get("Authorization")).To(Equal(""))
		})
	})
})

func makeRequest(method, path string) (*http.Request, error) {
	return http.NewRequest(method, path, nil)
}