// Command tiger-mcp is an MCP server exposing read-only Tiger Open API market
// data and account access to MCP clients over stdio.
package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/config"
	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
	"github.com/yogyrahmawan/tiger-mcp/internal/tools"
)

func main() {
	logger := log.New(os.Stderr, "tiger-mcp: ", log.LstdFlags)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("config error: %v", err)
	}

	tigerClient, err := tiger.NewClient(cfg)
	if err != nil {
		logger.Fatalf("tiger client error: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "tiger-mcp",
		Version: "0.0.1",
	}, nil)

	tools.RegisterGetQuote(server, tigerClient)
	tools.RegisterGetMarketStatus(server, tigerClient)
	tools.RegisterGetKline(server, tigerClient)
	tools.RegisterGetDepth(server, tigerClient)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Fatalf("server exited with error: %v", err)
	}
}
