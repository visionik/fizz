package format

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// TableFormatter formats output as a table
type TableFormatter struct {
	Writer io.Writer
}

// Format outputs data as a table
func (f *TableFormatter) Format(data interface{}) error {
	if data == nil {
		return nil
	}

	// Handle slices
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			fmt.Fprintln(f.Writer, "No results found")
			return nil
		}
		return f.formatSlice(v)
	}

	// Handle single objects
	return f.formatSingle(data)
}

func (f *TableFormatter) formatSlice(v reflect.Value) error {
	if v.Len() == 0 {
		return nil
	}

	// Get headers from first element
	first := v.Index(0)
	headers, err := f.getHeaders(first)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(f.Writer)
	// Convert []string to []any
	headerAny := make([]any, len(headers))
	for i, h := range headers {
		headerAny[i] = h
	}
	table.Header(headerAny...)

	// Add rows
	for i := 0; i < v.Len(); i++ {
		row, err := f.getRow(v.Index(i))
		if err != nil {
			return err
		}
		table.Append(row)
	}

	return table.Render()
}

func (f *TableFormatter) formatSingle(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		fmt.Fprintf(f.Writer, "%v\n", data)
		return nil
	}

	table := tablewriter.NewWriter(f.Writer)
	table.Header("Field", "Value")

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		fieldValue := fmt.Sprintf("%v", v.Field(i).Interface())

		// Add color for certain fields
		if strings.Contains(strings.ToLower(fieldName), "status") {
			fieldValue = f.colorizeStatus(fieldValue)
		}

		table.Append(fieldName, fieldValue)
	}

	return table.Render()
}

func (f *TableFormatter) getHeaders(v reflect.Value) ([]string, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return []string{"Value"}, nil
	}

	t := v.Type()
	headers := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		headers = append(headers, field.Name)
	}
	return headers, nil
}

func (f *TableFormatter) getRow(v reflect.Value) ([]string, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return []string{fmt.Sprintf("%v", v.Interface())}, nil
	}

	t := v.Type()
	row := make([]string, 0, t.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		value := fmt.Sprintf("%v", v.Field(i).Interface())

		// Add color for status fields
		if strings.Contains(strings.ToLower(field.Name), "status") {
			value = f.colorizeStatus(value)
		}

		row = append(row, value)
	}
	return row, nil
}

func (f *TableFormatter) colorizeStatus(status string) string {
	// Check if colors are disabled
	if color.NoColor {
		return status
	}

	lower := strings.ToLower(status)
	switch {
	case strings.Contains(lower, "open"), strings.Contains(lower, "active"):
		return color.GreenString(status)
	case strings.Contains(lower, "closed"), strings.Contains(lower, "completed"):
		return color.BlueString(status)
	case strings.Contains(lower, "error"), strings.Contains(lower, "failed"):
		return color.RedString(status)
	case strings.Contains(lower, "pending"), strings.Contains(lower, "waiting"):
		return color.YellowString(status)
	default:
		return status
	}
}
