package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for fizz.

To load completions:

Bash:
  $ source <(fizz completion bash)
  # To load completions for each session, add to your ~/.bashrc:
  $ fizz completion bash > /etc/bash_completion.d/fizz

Zsh:
  # If shell completion is not already enabled:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session:
  $ fizz completion zsh > "${fpath[1]}/_fizz"
  $ source ~/.zshrc

Fish:
  $ fizz completion fish | source
  # To load completions for each session:
  $ fizz completion fish > ~/.config/fish/completions/fizz.fish`,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return rootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
