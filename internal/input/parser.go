package input

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ReadInput reads input from stdin if "-" is provided, from a file if path is provided,
// or returns nil if neither
func ReadInput(inputPath string) ([]byte, error) {
	if inputPath == "" {
		return nil, nil
	}

	var reader io.Reader
	if inputPath == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open input file: %w", err)
		}
		defer file.Close()
		reader = file
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return data, nil
}

// ParseJSON parses JSON data into the target structure
func ParseJSON(data []byte, target interface{}) error {
	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

// MergeJSON reads input and merges it with existing data
func MergeJSON(inputPath string, target interface{}) error {
	data, err := ReadInput(inputPath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return ParseJSON(data, target)
}
