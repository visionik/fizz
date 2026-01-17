package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
)

var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Identity operations",
	Long:  "Get information about your Fizzy identity and accounts",
}

var identityGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get your identity information",
	Long:  "Retrieve your Fizzy identity, including user information and accounts",
	Example: `  fizz identity get
  fizz identity get --format=json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		
		identity, err := client.Identity.Get(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get identity: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(identity)
	},
}

func init() {
	identityCmd.AddCommand(identityGetCmd)
	rootCmd.AddCommand(identityCmd)
}
