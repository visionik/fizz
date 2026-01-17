package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	Token   string
	Account string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	token := os.Getenv("FIZZY_TOKEN")
	account := os.Getenv("FIZZY_ACCOUNT")

	if token == "" {
		return nil, fmt.Errorf(`FIZZY_TOKEN environment variable is not set

Setup instructions:
  1. Get your API token from https://fizzy.do/settings/tokens
  2. Set the environment variable:
     export FIZZY_TOKEN="your-token-here"
  3. Set your account ID:
     export FIZZY_ACCOUNT="your-account-id"

You can add these to your ~/.zshrc or ~/.bashrc to make them permanent.`)
	}

	if account == "" {
		return nil, fmt.Errorf(`FIZZY_ACCOUNT environment variable is not set

Setup instructions:
  1. Find your account ID at https://fizzy.do/settings/account
  2. Set the environment variable:
     export FIZZY_ACCOUNT="your-account-id"

You can add this to your ~/.zshrc or ~/.bashrc to make it permanent.`)
	}

	return &Config{
		Token:   token,
		Account: account,
	}, nil
}
