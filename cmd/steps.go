package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
	"github.com/visionik/libfizz-go/fizzy"
)

var stepsCmd = &cobra.Command{
	Use:   "steps",
	Short: "Checklist step management",
	Long:  "List, create, update, and delete checklist steps on cards",
}

var stepsListCmd = &cobra.Command{
	Use:   "list <card-id-or-number>",
	Short: "List steps on a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		steps, err := client.Steps.List(cmd.Context(), cardID)
		if err != nil {
			return fmt.Errorf("failed to list steps: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(steps)
	},
}

var stepsGetCmd = &cobra.Command{
	Use:   "get <card-id-or-number> <step-id>",
	Short: "Get a step",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		stepID := args[1]
		step, err := client.Steps.Get(cmd.Context(), cardID, stepID)
		if err != nil {
			return fmt.Errorf("failed to get step: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(step)
	},
}

var stepsCreateCmd = &cobra.Command{
	Use:   "create <card-id-or-number>",
	Short: "Create a checklist step",
	Args:  cobra.ExactArgs(1),
	Example: `  fizz steps create 123 --content="Review code"
  fizz steps create 123 --content="Deploy" --completed=true`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		content, _ := cmd.Flags().GetString("content")
		completed, _ := cmd.Flags().GetBool("completed")

		if content == "" {
			return fmt.Errorf("--content is required")
		}

	opts := &fizzy.StepCreateOptions{
		Content: content,
	}
	if completed {
		opts.Completed = &completed
	}

	step, err := client.Steps.Create(cmd.Context(), cardID, opts)
		if err != nil {
			return fmt.Errorf("failed to create step: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(step)
	},
}

var stepsUpdateCmd = &cobra.Command{
	Use:   "update <card-id-or-number> <step-id>",
	Short: "Update a step",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		stepID := args[1]
		content, _ := cmd.Flags().GetString("content")
		completed, _ := cmd.Flags().GetBool("completed")

		req := &fizzy.StepUpdateOptions{}
		if content != "" {
			req.Content = &content
		}
		if cmd.Flags().Changed("completed") {
			req.Completed = &completed
		}

		step, err := client.Steps.Update(cmd.Context(), cardID, stepID, req)
		if err != nil {
			return fmt.Errorf("failed to update step: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(step)
	},
}

var stepsDeleteCmd = &cobra.Command{
	Use:   "delete <card-id-or-number> <step-id>",
	Short: "Delete a step",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		cardID, err := client.ResolveCardID(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		stepID := args[1]
		err = client.Steps.Delete(cmd.Context(), cardID, stepID)
		if err != nil {
			return fmt.Errorf("failed to delete step: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Step deleted successfully\n")
		return nil
	},
}

func init() {
	stepsCreateCmd.Flags().String("content", "", "Step content (required)")
	stepsCreateCmd.Flags().Bool("completed", false, "Mark as completed")
	stepsUpdateCmd.Flags().String("content", "", "New step content")
	stepsUpdateCmd.Flags().Bool("completed", false, "Mark as completed")

	stepsCmd.AddCommand(stepsListCmd)
	stepsCmd.AddCommand(stepsGetCmd)
	stepsCmd.AddCommand(stepsCreateCmd)
	stepsCmd.AddCommand(stepsUpdateCmd)
	stepsCmd.AddCommand(stepsDeleteCmd)
	rootCmd.AddCommand(stepsCmd)
}
