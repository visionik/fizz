package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
	"github.com/visionik/libfizz-go/fizzy"
)

var cardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "Card operations",
	Long:  "Create, read, update, delete, and manage cards",
}

// List cards
var cardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cards",
	Example: `  fizz cards list
  fizz cards list --board=03fbhiu9dgjo0viyrlya1x03a
  fizz cards list --status=open --limit=20`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		limit, _ := cmd.Flags().GetInt("limit")
		boardID, _ := cmd.Flags().GetString("board")
		status, _ := cmd.Flags().GetString("status")

		// Build list options
		opts := &fizzy.CardListOptions{}
		if boardID != "" {
			opts.BoardID = boardID
		}
		if status != "" {
			opts.Status = status
		}

		cards, err := client.Cards.ListAll(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("failed to list cards: %w", err)
		}

		// Apply limit if specified
		if limit > 0 && len(cards) > limit {
			cards = cards[:limit]
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		// Use compact display for table format, full data for JSON/YAML
		if GetFormat() == "table" {
			return formatter.Format(format.ToCardDisplaySlice(cards))
		}
		return formatter.Format(cards)
	},
}

// Get card
var cardsGetCmd = &cobra.Command{
	Use:   "get <card-id-or-number>",
	Short: "Get a card",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz cards get 123
  fizz cards get abc-123-def --format=json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		
		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		card, err := client.Cards.Get(cmd.Context(), cardID)
		if err != nil {
			return fmt.Errorf("failed to get card: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		// Use detail display for table format, full data for JSON/YAML
		if GetFormat() == "table" {
			return formatter.Format(format.ToCardDetailDisplay(*card))
		}
		return formatter.Format(card)
	},
}

// Create card
var cardsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new card",
	Example: `  fizz cards create --board=03fbhiu9dgjo0viyrlya1x03a --title="Bug fix"
  fizz cards create --board=03fbhiu9dgjo0viyrlya1x03a --title="Feature" --body="Description"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		boardID, _ := cmd.Flags().GetString("board")
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")

		if boardID == "" {
			return fmt.Errorf("--board is required")
		}
		if title == "" {
			return fmt.Errorf("--title is required")
		}

		opts := &fizzy.CardCreateOptions{
			BoardID: boardID,
			Title:   title,
		}
		if body != "" {
			opts.Body = &body
		}

		card, err := client.Cards.Create(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("failed to create card: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(card)
	},
}

// Update card
var cardsUpdateCmd = &cobra.Command{
	Use:   "update <card-id-or-number>",
	Short: "Update a card",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz cards update 123 --title="Updated title"
  fizz cards update 123 --body="New description"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")

		opts := &fizzy.CardUpdateOptions{}
		if title != "" {
			opts.Title = &title
		}
		if body != "" {
			opts.Body = &body
		}

		card, err := client.Cards.Update(cmd.Context(), cardID, opts)
		if err != nil {
			return fmt.Errorf("failed to update card: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(card)
	},
}

// Delete card
var cardsDeleteCmd = &cobra.Command{
	Use:   "delete <card-id-or-number>",
	Short: "Delete a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		err = client.Cards.Delete(cmd.Context(), cardID)
		if err != nil {
			return fmt.Errorf("failed to delete card: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Card %d deleted successfully\n", cardID)
		return nil
	},
}

// Close card
var cardsCloseCmd = &cobra.Command{
	Use:   "close <card-id-or-number>",
	Short: "Close a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

	err = client.Cards.Close(cmd.Context(), cardID)
	if err != nil {
		return fmt.Errorf("failed to close card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Card %s closed successfully\n", cardID)
	return nil
	},
}

// Reopen card
var cardsReopenCmd = &cobra.Command{
	Use:   "reopen <card-id-or-number>",
	Short: "Reopen a closed card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

	err = client.Cards.Reopen(cmd.Context(), cardID)
	if err != nil {
		return fmt.Errorf("failed to reopen card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Card %s reopened successfully\n", cardID)
	return nil
	},
}

// Postpone card
var cardsPostponeCmd = &cobra.Command{
	Use:   "postpone <card-id-or-number>",
	Short: "Postpone a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

	err = client.Cards.Postpone(cmd.Context(), cardID)
	if err != nil {
		return fmt.Errorf("failed to postpone card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Card %s postponed successfully\n", cardID)
	return nil
	},
}

// Triage card
var cardsTriageCmd = &cobra.Command{
	Use:   "triage <card-id-or-number>",
	Short: "Triage a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

	err = client.Cards.Triage(cmd.Context(), cardID)
	if err != nil {
		return fmt.Errorf("failed to triage card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Card %s triaged successfully\n", cardID)
	return nil
	},
}

// Assign card
var cardsAssignCmd = &cobra.Command{
	Use:   "assign <card-id-or-number> <user-id>",
	Short: "Assign a card to a user",
	Args:  cobra.ExactArgs(2),
	Example: `  fizz cards assign 123 user-456`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		userID := args[1]

	err = client.Cards.Assign(cmd.Context(), cardID, userID)
	if err != nil {
		return fmt.Errorf("failed to assign card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Card %s assigned to %s\n", cardID, userID)
	return nil
	},
}

// Tag card
var cardsTagCmd = &cobra.Command{
	Use:   "tag <card-id-or-number> <tag-name>",
	Short: "Add a tag to a card",
	Args:  cobra.ExactArgs(2),
	Example: `  fizz cards tag 123 bug`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		tagName := args[1]

	err = client.Cards.Tag(cmd.Context(), cardID, tagName)
	if err != nil {
		return fmt.Errorf("failed to tag card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Card %s tagged with '%s'\n", cardID, tagName)
	return nil
	},
}

// Move card
var cardsMoveCmd = &cobra.Command{
	Use:   "move <card-id-or-number>",
	Short: "Move a card to a different column",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz cards move 123 --column=456`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		columnID, _ := cmd.Flags().GetInt("column")
		if columnID == 0 {
			return fmt.Errorf("--column is required")
		}

	err = client.Cards.MoveToColumn(cmd.Context(), cardID, strconv.Itoa(columnID))
	if err != nil {
		return fmt.Errorf("failed to move card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Card %s moved to column %d\n", cardID, columnID)
	return nil
	},
}

// Watch card
var cardsWatchCmd = &cobra.Command{
	Use:   "watch <card-id-or-number>",
	Short: "Watch a card for updates",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

	err = client.Cards.Watch(cmd.Context(), cardID)
	if err != nil {
		return fmt.Errorf("failed to watch card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Now watching card %s\n", cardID)
	return nil
	},
}

// Unwatch card
var cardsUnwatchCmd = &cobra.Command{
	Use:   "unwatch <card-id-or-number>",
	Short: "Stop watching a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

	err = client.Cards.Unwatch(cmd.Context(), cardID)
	if err != nil {
		return fmt.Errorf("failed to unwatch card: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Stopped watching card %s\n", cardID)
	return nil
	},
}

// Golden card
var cardsGoldenCmd = &cobra.Command{
	Use:   "golden <card-id-or-number>",
	Short: "Mark a card as golden",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

	err = client.Cards.MarkGolden(cmd.Context(), cardID)
	if err != nil {
		return fmt.Errorf("failed to mark card as golden: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Card %s marked as golden\n", cardID)
	return nil
	},
}

// Ungolden card
var cardsUngoldenCmd = &cobra.Command{
	Use:   "ungolden <card-id-or-number>",
	Short: "Remove golden status from a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

	err = client.Cards.UnmarkGolden(cmd.Context(), cardID)
	if err != nil {
		return fmt.Errorf("failed to remove golden status: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Removed golden status from card %s\n", cardID)
	return nil
	},
}

func init() {
	// List flags
	cardsListCmd.Flags().String("board", "", "Filter by board ID")
	cardsListCmd.Flags().String("status", "", "Filter by status")
	cardsListCmd.Flags().Int("limit", 0, "Limit number of results (0 = all)")

	// Create flags
	cardsCreateCmd.Flags().String("board", "", "Board ID (required)")
	cardsCreateCmd.Flags().String("title", "", "Card title (required)")
	cardsCreateCmd.Flags().String("body", "", "Card body/description")

	// Update flags
	cardsUpdateCmd.Flags().String("title", "", "New card title")
	cardsUpdateCmd.Flags().String("body", "", "New card body")

	// Move flags
	cardsMoveCmd.Flags().Int("column", 0, "Target column ID (required)")

	// Add all subcommands
	cardsCmd.AddCommand(cardsListCmd)
	cardsCmd.AddCommand(cardsGetCmd)
	cardsCmd.AddCommand(cardsCreateCmd)
	cardsCmd.AddCommand(cardsUpdateCmd)
	cardsCmd.AddCommand(cardsDeleteCmd)
	cardsCmd.AddCommand(cardsCloseCmd)
	cardsCmd.AddCommand(cardsReopenCmd)
	cardsCmd.AddCommand(cardsPostponeCmd)
	cardsCmd.AddCommand(cardsTriageCmd)
	cardsCmd.AddCommand(cardsAssignCmd)
	cardsCmd.AddCommand(cardsTagCmd)
	cardsCmd.AddCommand(cardsMoveCmd)
	cardsCmd.AddCommand(cardsWatchCmd)
	cardsCmd.AddCommand(cardsUnwatchCmd)
	cardsCmd.AddCommand(cardsGoldenCmd)
	cardsCmd.AddCommand(cardsUngoldenCmd)

	rootCmd.AddCommand(cardsCmd)
}
