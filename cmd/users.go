package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User operations",
	Long:  "List users in the account",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		users, err := client.Users.List(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(users)
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
	rootCmd.AddCommand(usersCmd)
}
