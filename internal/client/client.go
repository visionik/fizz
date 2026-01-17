package client

import (
	"fmt"
	"log"

	"github.com/visionik/fizz/internal/config"
	"github.com/visionik/libfizz-go/fizzy"
)

// Client wraps the libfizz-go client
type Client struct {
	*fizzy.Client
	Debug bool
}

// New creates a new Fizzy client from configuration
func New(cfg *config.Config, debug bool) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	client := fizzy.NewClient(cfg.Token, cfg.Account)

	if debug {
		log.Println("Debug mode enabled")
		log.Printf("Account ID: %s", cfg.Account)
	}

	return &Client{
		Client: client,
		Debug:  debug,
	}, nil
}
