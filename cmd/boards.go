package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
	"github.com/visionik/libfizz-go/fizzy"
)

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Board management",
	Long:  "List, create, update, and delete boards",
}

var boardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all boards",
	Long:  "List all boards in your account",
	Example: `  fizz boards list
  fizz boards list --format=json
  fizz boards list --limit=10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		limit, _ := cmd.Flags().GetInt("limit")

		boards, err := client.Boards.List(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to list boards: %w", err)
		}

		if limit > 0 && len(boards) > limit {
			boards = boards[:limit]
		}

		if err != nil {
			return fmt.Errorf("failed to list boards: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(boards)
	},
}

var boardsGetCmd = &cobra.Command{
	Use:   "get <board-id>",
	Short: "Get a board by ID",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz boards get 123
  fizz boards get 123 --format=json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		boardID := args[0]

		board, err := client.Boards.Get(cmd.Context(), boardID)
		if err != nil {
			return fmt.Errorf("failed to get board: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(board)
	},
}

var boardsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new board",
	Long:  "Create a new board with the specified name and optional description",
	Example: `  fizz boards create --name="My Board"
  fizz boards create --name="My Board" --description="Board description"
  echo '{"name":"My Board","description":"Desc"}' | fizz boards create --input=-`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		opts := &fizzy.BoardCreateOptions{
			Name: name,
		}
		if description != "" {
			opts.Description = &description
		}

		board, err := client.Boards.Create(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("failed to create board: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(board)
	},
}

var boardsUpdateCmd = &cobra.Command{
	Use:   "update <board-id>",
	Short: "Update a board",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz boards update 123 --name="Updated Name"
  fizz boards update 123 --description="Updated description"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		boardID := args[0]

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		opts := &fizzy.BoardUpdateOptions{}
		if name != "" {
			opts.Name = &name
		}
		if description != "" {
			opts.Description = &description
		}

		board, err := client.Boards.Update(cmd.Context(), boardID, opts)
		if err != nil {
			return fmt.Errorf("failed to update board: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(board)
	},
}

var boardsDeleteCmd = &cobra.Command{
	Use:   "delete <board-id>",
	Short: "Delete a board",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz boards delete 123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		boardID := args[0]

		err := client.Boards.Delete(cmd.Context(), boardID)
		if err != nil {
			return fmt.Errorf("failed to delete board: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Board %s deleted successfully\n", boardID)
		return nil
	},
}

func init() {
	boardsListCmd.Flags().Int("limit", 0, "Limit number of results (0 = all)")
	boardsCreateCmd.Flags().String("name", "", "Board name (required)")
	boardsCreateCmd.Flags().String("description", "", "Board description")
	boardsCreateCmd.Flags().String("input", "", "Read from file or stdin (-)")
	boardsUpdateCmd.Flags().String("name", "", "New board name")
	boardsUpdateCmd.Flags().String("description", "", "New board description")

	boardsCmd.AddCommand(boardsListCmd)
	boardsCmd.AddCommand(boardsGetCmd)
	boardsCmd.AddCommand(boardsCreateCmd)
	boardsCmd.AddCommand(boardsUpdateCmd)
	boardsCmd.AddCommand(boardsDeleteCmd)
	rootCmd.AddCommand(boardsCmd)
}
