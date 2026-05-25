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
			client := api.NewClient("https://open-notebook.darimuri.me", auth.NewMiddleware("test-key"))
			Expect(client.BaseURL()).To(Equal("https://open-notebook.darimuri.me"))
		})
	})

	Describe("NewRequest", func() {
		It("sets auth header on request", func() {
			client := api.NewClient("https://open-notebook.darimuri.me", auth.NewMiddleware("test-key"))
			req, err := client.NewRequest("GET", "/api/notebooks", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(req.Header.Get("Authorization")).To(Equal("Bearer test-key"))
		})

		It("sets Content-Type header", func() {
			client := api.NewClient("https://open-notebook.darimuri.me", auth.NewMiddleware(""))
			req, err := client.NewRequest("POST", "/api/notebooks", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(req.Header.Get("Content-Type")).To(Equal("application/json"))
		})

		It("sets Accept header", func() {
			client := api.NewClient("https://open-notebook.darimuri.me", auth.NewMiddleware(""))
			req, err := client.NewRequest("GET", "/api/sources", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(req.Header.Get("Accept")).To(Equal("application/json"))
		})
	})

	Describe("Get sources with embedded field", func() {
		It("returns embedded field for sources", func() {
			client := api.NewClient("https://open-notebook.darimuri.me", auth.NewMiddleware(""))
			var sources []map[string]interface{}
			err := client.Get("/api/sources?limit=10", &sources)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(sources)).To(BeNumerically(">", 0))

			// Check that embedded field is present
			source := sources[0]
			Expect(source).To(HaveKey("embedded"), "embedded field should be present in response")
		})

		It("source:wc2mquohoo12nqsl0vov should have embedded=true", func() {
			client := api.NewClient("https://open-notebook.darimuri.me", auth.NewMiddleware(""))
			var sources []map[string]interface{}
			err := client.Get("/api/sources?limit=100", &sources)
			Expect(err).NotTo(HaveOccurred())

			var foundSource map[string]interface{}
			for _, s := range sources {
				if s["id"] == "source:wc2mquohoo12nqsl0vov" {
					foundSource = s
					break
				}
			}
			Expect(foundSource).NotTo(BeNil(), "source:wc2mquohoo12nqsl0vov should be found")
			Expect(foundSource).To(HaveKey("embedded"), "embedded field should be present")
			Expect(foundSource["embedded"]).To(BeEquivalentTo(true), "source:wc2mquohoo12nqsl0vov should have embedded=true")
		})

		It("source:d89jd8j1uva4mk4ohh1k should have embedded=nil or false", func() {
			client := api.NewClient("https://open-notebook.darimuri.me", auth.NewMiddleware(""))
			var sources []map[string]interface{}
			err := client.Get("/api/sources?limit=100", &sources)
			Expect(err).NotTo(HaveOccurred())

			var foundSource map[string]interface{}
			for _, s := range sources {
				if s["id"] == "source:d89jd8j1uva4mk4ohh1k" {
					foundSource = s
					break
				}
			}
			Expect(foundSource).NotTo(BeNil(), "source:d89jd8j1uva4mk4ohh1k should be found")
			Expect(foundSource).To(HaveKey("embedded"), "embedded field should be present")
			embeddedVal := foundSource["embedded"]
			Expect(embeddedVal == nil || embeddedVal == false).To(BeTrue(), "source:d89jd8j1uva4mk4ohh1k should have embedded=nil or false, got: %v", embeddedVal)
		})
	})
})