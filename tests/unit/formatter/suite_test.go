package formatter_test

import (
	"bytes"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/darimuri/open-notebook-cli/internal/formatter"
)

func TestFormatter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Formatter Suite")
}

var _ = Describe("Formatter", func() {
	Describe("JSON format", func() {
		It("outputs valid JSON", func() {
			f := formatter.New("json")
			data := map[string]any{"name": "test", "value": 123}

			buf := &bytes.Buffer{}
			err := f.Format(buf, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring(`"name"`))
			Expect(buf.String()).To(ContainSubstring(`"test"`))
		})

		It("formats slices as JSON array", func() {
			f := formatter.New("json")
			data := []map[string]string{
				{"id": "1", "name": "first"},
				{"id": "2", "name": "second"},
			}

			buf := &bytes.Buffer{}
			err := f.Format(buf, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring(`"id"`))
			Expect(buf.String()).To(ContainSubstring(`"first"`))
		})
	})

	Describe("Table format", func() {
		It("outputs tab-separated format", func() {
			f := formatter.New("table")
			data := []map[string]string{
				{"id": "1", "name": "first"},
				{"id": "2", "name": "second"},
			}

			buf := &bytes.Buffer{}
			err := f.Format(buf, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("id"))
			Expect(buf.String()).To(ContainSubstring("name"))
			Expect(buf.String()).To(ContainSubstring("1"))
			Expect(buf.String()).To(ContainSubstring("first"))
		})

		It("handles empty slice", func() {
			f := formatter.New("table")
			data := []map[string]string{}

			buf := &bytes.Buffer{}
			err := f.Format(buf, data)

			Expect(err).NotTo(HaveOccurred())
			Expect(buf.String()).To(ContainSubstring("(empty)"))
		})

		It("returns error for non-slice data", func() {
			f := formatter.New("table")
			data := map[string]string{"name": "test"}

			buf := &bytes.Buffer{}
			err := f.Format(buf, data)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("table format requires"))
		})
	})
})