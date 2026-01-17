package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
	"github.com/visionik/libfizz-go/fizzy"
)

var columnsCmd = &cobra.Command{
	Use:   "columns",
	Short: "Column management",
	Long:  "List, create, update, and delete columns",
}

var columnsListCmd = &cobra.Command{
	Use:   "list <board-id>",
	Short: "List columns in a board",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		boardID := args[0]

		columns, err := client.Columns.List(cmd.Context(), boardID)
		if err != nil {
			return fmt.Errorf("failed to list columns: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(columns)
	},
}

var columnsGetCmd = &cobra.Command{
	Use:   "get <board-id> <column-id>",
	Short: "Get a column",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		boardID := args[0]
		columnID := args[1]

		column, err := client.Columns.Get(cmd.Context(), boardID, columnID)
		if err != nil {
			return fmt.Errorf("failed to get column: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(column)
	},
}

var columnsCreateCmd = &cobra.Command{
	Use:   "create <board-id>",
	Short: "Create a column",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz columns create 123 --name="In Progress"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		boardID := args[0]

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}

	opts := &fizzy.ColumnCreateOptions{
		Name: name,
	}

	column, err := client.Columns.Create(cmd.Context(), boardID, opts)
		if err != nil {
			return fmt.Errorf("failed to create column: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(column)
	},
}

var columnsUpdateCmd = &cobra.Command{
	Use:   "update <board-id> <column-id>",
	Short: "Update a column",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		boardID := args[0]
		columnID := args[1]

		name, _ := cmd.Flags().GetString("name")
		position, _ := cmd.Flags().GetInt("position")

	opts := &fizzy.ColumnUpdateOptions{}
	if name != "" {
		opts.Name = &name
	}
	if position > 0 {
		opts.Position = &position
	}

	column, err := client.Columns.Update(cmd.Context(), boardID, columnID, opts)
		if err != nil {
			return fmt.Errorf("failed to update column: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(column)
	},
}

var columnsDeleteCmd = &cobra.Command{
	Use:   "delete <board-id> <column-id>",
	Short: "Delete a column",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		boardID := args[0]
		columnID := args[1]

		err := client.Columns.Delete(cmd.Context(), boardID, columnID)
		if err != nil {
			return fmt.Errorf("failed to delete column: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Column deleted successfully\n")
		return nil
	},
}

func init() {
	columnsCreateCmd.Flags().String("name", "", "Column name (required)")
	columnsUpdateCmd.Flags().String("name", "", "New column name")
	columnsUpdateCmd.Flags().Int("position", 0, "New position")

	columnsCmd.AddCommand(columnsListCmd)
	columnsCmd.AddCommand(columnsGetCmd)
	columnsCmd.AddCommand(columnsCreateCmd)
	columnsCmd.AddCommand(columnsUpdateCmd)
	columnsCmd.AddCommand(columnsDeleteCmd)
	rootCmd.AddCommand(columnsCmd)
}
