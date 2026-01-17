package client

import (
	"context"
)


// ResolveCardID takes a card identifier (number or UUID) and returns the card ID
// The Fizzy API uses string IDs, so we just return the input as-is
func (c *Client) ResolveCardID(ctx context.Context, input string) (string, error) {
	// Just return the input - API accepts both IDs and numbers as strings
	return input, nil
}
