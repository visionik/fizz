package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uploadsCmd = &cobra.Command{
	Use:   "uploads",
	Short: "File upload operations",
	Long:  "Upload files to Fizzy",
}

var uploadsCreateCmd = &cobra.Command{
	Use:   "create <file-path>",
	Short: "Upload a file",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz uploads create ./image.png
  fizz uploads create /path/to/document.pdf`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		filePath := args[0]

		// Upload file (content type will be auto-detected)
		url, err := client.Uploads.UploadFile(cmd.Context(), filePath, "")
		if err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "File uploaded successfully\nURL: %s\n", url)
		return nil
	},
}

func init() {
	uploadsCmd.AddCommand(uploadsCreateCmd)
	rootCmd.AddCommand(uploadsCmd)
}
