package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
	"github.com/visionik/libfizz-go/fizzy"
)

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Comment management",
	Long:  "List, create, update, and delete comments on cards",
}

var commentsListCmd = &cobra.Command{
	Use:   "list <card-id-or-number>",
	Short: "List comments on a card",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz comments list 123
  fizz comments list abc-def --format=json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		comments, err := client.Comments.List(cmd.Context(), cardID)
		if err != nil {
			return fmt.Errorf("failed to list comments: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(comments)
	},
}

var commentsCreateCmd = &cobra.Command{
	Use:   "create <card-id-or-number>",
	Short: "Create a comment on a card",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz comments create 123 --body="Great work!"
  echo "Long comment text" | fizz comments create 123 --body=-`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		body, _ := cmd.Flags().GetString("body")
		if body == "" {
			return fmt.Errorf("--body is required")
		}

		req := &fizzy.CommentCreateOptions{
			Body: body,
		}

		comment, err := client.Comments.Create(cmd.Context(), cardID, req)
		if err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(comment)
	},
}

var commentsUpdateCmd = &cobra.Command{
	Use:   "update <card-id-or-number> <comment-id>",
	Short: "Update a comment",
	Args:  cobra.ExactArgs(2),
	Example: `  fizz comments update 123 456 --body="Updated comment"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		commentID := args[1]
		body, _ := cmd.Flags().GetString("body")
		if body == "" {
			return fmt.Errorf("--body is required")
		}

		req := &fizzy.CommentUpdateOptions{
			Body: body,
		}

		comment, err := client.Comments.Update(cmd.Context(), cardID, commentID, req)
		if err != nil {
			return fmt.Errorf("failed to update comment: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(comment)
	},
}

var commentsDeleteCmd = &cobra.Command{
	Use:   "delete <card-id-or-number> <comment-id>",
	Short: "Delete a comment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		commentID := args[1]

		err = client.Comments.Delete(cmd.Context(), cardID, commentID)
		if err != nil {
			return fmt.Errorf("failed to delete comment: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Comment %s deleted successfully\n", commentID)
		return nil
	},
}

func init() {
	commentsCreateCmd.Flags().String("body", "", "Comment body (required)")
	commentsUpdateCmd.Flags().String("body", "", "New comment body (required)")

	commentsCmd.AddCommand(commentsListCmd)
	commentsCmd.AddCommand(commentsCreateCmd)
	commentsCmd.AddCommand(commentsUpdateCmd)
	commentsCmd.AddCommand(commentsDeleteCmd)
	rootCmd.AddCommand(commentsCmd)
}
