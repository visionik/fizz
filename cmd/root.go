package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/aihelp"
	"github.com/visionik/fizz/internal/client"
	"github.com/visionik/fizz/internal/config"
)

var (
	formatFlag string
	debugFlag  bool
	globalClient *client.Client
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "fizz",
	Short: "A CLI for Fizzy.do - Kanban project management",
	Long: `fizz is a command-line interface for Fizzy.do (https://fizzy.do).

It provides a simple, consistent way to interact with all Fizzy API features
through the terminal. Use 'fizz <noun> <verb>' to perform operations.

Authentication:
  Set FIZZY_TOKEN and FIZZY_ACCOUNT environment variables.

Examples:
  fizz boards list
  fizz cards create --board=123 --title="Bug fix"
  fizz identity get --format=json

For more help on any command, use: fizz <command> --help
Need AI-powered help? Use the --ai-help flag on any command.`,
	PersistentPreRunE: initClient,
}

// Execute runs the root command
func Execute() {
	// Check for --ai-help flag before processing anything else
	for _, arg := range os.Args[1:] {
		if arg == "--ai-help" {
			fmt.Fprint(os.Stdout, getAIHelp())
			os.Exit(0)
		}
	}
	
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "table", "Output format (table, json, yaml)")
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug output")
	
	// Add --ai-help flag
	rootCmd.PersistentFlags().Bool("ai-help", false, "Show AI-powered help for this command")
	
	// Set help template to mention --ai-help
	rootCmd.SetHelpTemplate(rootCmd.HelpTemplate() + "\nAI/LLMs SHOULD do a \"fizz --ai-help\"\n")
}

// initClient initializes the Fizzy client
func initClient(cmd *cobra.Command, args []string) error {
	// Skip client initialization for certain commands
	if cmd.Name() == "help" || cmd.Name() == "completion" {
		return nil
	}

	// Check for --ai-help flag
	aiHelp, _ := cmd.Flags().GetBool("ai-help")
	if aiHelp {
		// Import is needed at top of file
		fmt.Fprint(os.Stdout, getAIHelp())
		os.Exit(0)
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {
		return err
	}

	globalClient, err = client.New(cfg, debugFlag)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	return nil
}

// GetClient returns the global client instance
func GetClient() *client.Client {
	return globalClient
}

// GetFormat returns the format flag value
func GetFormat() string {
	return formatFlag
}

// getAIHelp returns the AI help content
func getAIHelp() string {
	return aihelp.Content
}
