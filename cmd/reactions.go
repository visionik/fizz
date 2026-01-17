package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
	"github.com/visionik/libfizz-go/fizzy"
)

var reactionsCmd = &cobra.Command{
	Use:   "reactions",
	Short: "Reaction management",
	Long:  "List, create, and delete reactions on comments",
}

var reactionsListCmd = &cobra.Command{
	Use:   "list <card-id-or-number> <comment-id>",
	Short: "List reactions on a comment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		commentID := args[1]

		reactions, err := client.Reactions.List(cmd.Context(), cardID, commentID)
		if err != nil {
			return fmt.Errorf("failed to list reactions: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(reactions)
	},
}

var reactionsCreateCmd = &cobra.Command{
	Use:   "create <card-id-or-number> <comment-id>",
	Short: "Add a reaction to a comment",
	Args:  cobra.ExactArgs(2),
	Example: `  fizz reactions create 123 456 --emoji="👍"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		commentID := args[1]
		emoji, _ := cmd.Flags().GetString("emoji")
		if emoji == "" {
			return fmt.Errorf("--emoji is required")
		}

	opts := &fizzy.ReactionCreateOptions{
		Content: emoji,
	}

	reaction, err := client.Reactions.Create(cmd.Context(), cardID, commentID, opts)
		if err != nil {
			return fmt.Errorf("failed to create reaction: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(reaction)
	},
}

var reactionsDeleteCmd = &cobra.Command{
	Use:   "delete <card-id-or-number> <comment-id> <reaction-id>",
	Short: "Delete a reaction",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		commentID := args[1]
		reactionID := args[2]

		err = client.Reactions.Delete(cmd.Context(), cardID, commentID, reactionID)
		if err != nil {
			return fmt.Errorf("failed to delete reaction: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Reaction deleted successfully\n")
		return nil
	},
}

func init() {
	reactionsCreateCmd.Flags().String("emoji", "", "Emoji reaction (required)")

	reactionsCmd.AddCommand(reactionsListCmd)
	reactionsCmd.AddCommand(reactionsCreateCmd)
	reactionsCmd.AddCommand(reactionsDeleteCmd)
	rootCmd.AddCommand(reactionsCmd)
}
