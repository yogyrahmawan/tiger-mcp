// Command tiger-mcp is an MCP server exposing read-only Tiger Open API market
// data and account access to MCP clients over stdio.
package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	logger := log.New(os.Stderr, "tiger-mcp: ", log.LstdFlags)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "tiger-mcp",
		Version: "0.0.1",
	}, nil)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Fatalf("server exited with error: %v", err)
	}
}
