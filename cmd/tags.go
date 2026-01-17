package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
	"github.com/visionik/libfizz-go/fizzy"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Tag management",
	Long:  "List and create tags",
}

var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		tags, err := client.Tags.List(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to list tags: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(tags)
	},
}

var tagsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new tag",
	Example: `  fizz tags create --name="bug"
  fizz tags create --name="feature" --color="#00ff00"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

	opts := &fizzy.TagCreateOptions{
		Name: name,
	}
	if color != "" {
		opts.Color = &color
	}

	tag, err := client.Tags.Create(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("failed to create tag: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(tag)
	},
}

func init() {
	tagsCreateCmd.Flags().String("name", "", "Tag name (required)")
	tagsCreateCmd.Flags().String("color", "", "Tag color (hex)")

	tagsCmd.AddCommand(tagsListCmd)
	tagsCmd.AddCommand(tagsCreateCmd)
	rootCmd.AddCommand(tagsCmd)
}
