package format

import (
	"fmt"
	"io"
	"os"
)

// Formatter formats data for output
type Formatter interface {
	Format(data interface{}) error
}

// NewFormatter creates a formatter based on the format string
func NewFormatter(format string, writer io.Writer) (Formatter, error) {
	if writer == nil {
		writer = os.Stdout
	}

	switch format {
	case "table", "":
		return &TableFormatter{Writer: writer}, nil
	case "json":
		return &JSONFormatter{Writer: writer}, nil
	case "yaml":
		return &YAMLFormatter{Writer: writer}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s (supported: table, json, yaml)", format)
	}
}
