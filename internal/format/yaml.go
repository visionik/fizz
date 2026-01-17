package format

import (
	"io"

	"gopkg.in/yaml.v3"
)

// YAMLFormatter formats output as YAML
type YAMLFormatter struct {
	Writer io.Writer
}

// Format outputs data as YAML
func (f *YAMLFormatter) Format(data interface{}) error {
	encoder := yaml.NewEncoder(f.Writer)
	encoder.SetIndent(2)
	defer encoder.Close()
	return encoder.Encode(data)
}
