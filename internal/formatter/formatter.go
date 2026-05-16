package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Formatter struct {
	format string
}

func New(format string) *Formatter {
	return &Formatter{format: format}
}

func (f *Formatter) Format(w io.Writer, data interface{}) error {
	switch f.format {
	case "json":
		return f.formatJSON(w, data)
	case "table":
		return f.formatTable(w, data)
	default:
		return f.formatTable(w, data)
	}
}

func (f *Formatter) formatJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (f *Formatter) formatTable(w io.Writer, data interface{}) error {
	slice, ok := data.([]map[string]string)
	if !ok {
		return fmt.Errorf("table format requires []map[string]string")
	}

	if len(slice) == 0 {
		fmt.Fprintln(w, "(empty)")
		return nil
	}

	// Get headers
	var headers []string
	for key := range slice[0] {
		headers = append(headers, key)
	}

	// Print header
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print separator
	fmt.Fprintln(w, strings.Repeat("-", len(strings.Join(headers, "\t"))))

	// Print rows
	for _, row := range slice {
		var values []string
		for _, h := range headers {
			values = append(values, row[h])
		}
		fmt.Fprintln(w, strings.Join(values, "\t"))
	}

	return nil
}