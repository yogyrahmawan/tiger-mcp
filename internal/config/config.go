// Package config loads and validates Tiger Open API credentials from the
// environment.
package config

import (
	"fmt"
	"os"
)

// Config holds the Tiger Open API credentials needed to build a client.
type Config struct {
	TigerID    string
	PrivateKey string
	Account    string
}

// Load reads TIGER_ID, TIGER_PRIVATE_KEY, and TIGER_ACCOUNT from the
// environment. It fails fast, naming the specific missing variable, if any
// of them is empty.
func Load() (*Config, error) {
	tigerID := os.Getenv("TIGER_ID")
	if tigerID == "" {
		return nil, fmt.Errorf("config: TIGER_ID is required but not set")
	}

	privateKey := os.Getenv("TIGER_PRIVATE_KEY")
	if privateKey == "" {
		return nil, fmt.Errorf("config: TIGER_PRIVATE_KEY is required but not set")
	}

	account := os.Getenv("TIGER_ACCOUNT")
	if account == "" {
		return nil, fmt.Errorf("config: TIGER_ACCOUNT is required but not set")
	}

	return &Config{
		TigerID:    tigerID,
		PrivateKey: privateKey,
		Account:    account,
	}, nil
}
