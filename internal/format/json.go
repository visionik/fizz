package format

import (
	"encoding/json"
	"io"
)

// JSONFormatter formats output as JSON
type JSONFormatter struct {
	Writer io.Writer
}

// Format outputs data as JSON
func (f *JSONFormatter) Format(data interface{}) error {
	encoder := json.NewEncoder(f.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
