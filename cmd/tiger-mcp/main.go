// Command tiger-mcp is an MCP server exposing read-only Tiger Open API market
// data and account access to MCP clients over stdio.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/config"
	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
	"github.com/yogyrahmawan/tiger-mcp/internal/tools"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config error", "err", err)
		os.Exit(1)
	}

	tigerClient, err := tiger.NewClient(cfg)
	if err != nil {
		logger.Error("tiger client error", "err", err)
		os.Exit(1)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "tiger-mcp",
		Version: "0.0.1",
	}, nil)

	tools.RegisterGetQuote(server, tigerClient)
	tools.RegisterGetMarketStatus(server, tigerClient)
	tools.RegisterGetKline(server, tigerClient)
	tools.RegisterGetDepth(server, tigerClient)
	tools.RegisterGetAccountAssets(server, tigerClient)
	tools.RegisterGetPositions(server, tigerClient)
	tools.RegisterGetOrders(server, tigerClient)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
