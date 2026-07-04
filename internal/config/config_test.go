package config

import (
	"strings"
	"testing"
)

func TestLoad_AllPresent(t *testing.T) {
	t.Setenv("TIGER_ID", "tiger-id-1")
	t.Setenv("TIGER_PRIVATE_KEY", "-----BEGIN PRIVATE KEY-----\nabc\n-----END PRIVATE KEY-----")
	t.Setenv("TIGER_ACCOUNT", "account-1")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}
	if cfg.TigerID != "tiger-id-1" {
		t.Errorf("TigerID = %q, want %q", cfg.TigerID, "tiger-id-1")
	}
	if cfg.PrivateKey != "-----BEGIN PRIVATE KEY-----\nabc\n-----END PRIVATE KEY-----" {
		t.Errorf("PrivateKey mismatch: got %q", cfg.PrivateKey)
	}
	if cfg.Account != "account-1" {
		t.Errorf("Account = %q, want %q", cfg.Account, "account-1")
	}
}

func TestLoad_MissingVars(t *testing.T) {
	tests := []struct {
		name        string
		tigerID     string
		privateKey  string
		account     string
		wantErrText string
	}{
		{
			name:        "missing TIGER_ID",
			tigerID:     "",
			privateKey:  "key",
			account:     "account",
			wantErrText: "TIGER_ID",
		},
		{
			name:        "missing TIGER_PRIVATE_KEY",
			tigerID:     "id",
			privateKey:  "",
			account:     "account",
			wantErrText: "TIGER_PRIVATE_KEY",
		},
		{
			name:        "missing TIGER_ACCOUNT",
			tigerID:     "id",
			privateKey:  "key",
			account:     "",
			wantErrText: "TIGER_ACCOUNT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TIGER_ID", tt.tigerID)
			t.Setenv("TIGER_PRIVATE_KEY", tt.privateKey)
			t.Setenv("TIGER_ACCOUNT", tt.account)

			_, err := Load()
			if err == nil {
				t.Fatal("Load() returned nil error, want error naming the missing var")
			}
			if got := err.Error(); !strings.Contains(got, tt.wantErrText) {
				t.Errorf("Load() error = %q, want it to mention %q", got, tt.wantErrText)
			}
		})
	}
}
