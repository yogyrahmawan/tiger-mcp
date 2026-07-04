// Package tools defines the MCP tools exposed by this server, backed by the
// internal/tiger client wrapper.
package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

// GetQuoteInput is the input schema for the get_quote tool.
type GetQuoteInput struct {
	Symbols []string `json:"symbols" jsonschema:"ticker symbols to fetch real-time quotes for, e.g. AAPL, TSLA"`
}

// GetQuoteOutput is the output schema for the get_quote tool.
type GetQuoteOutput struct {
	Quotes []tiger.Quote `json:"quotes"`
}

// RegisterGetQuote adds the get_quote tool to server, backed by fetcher.
func RegisterGetQuote(server *mcp.Server, fetcher tiger.QuoteFetcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_quote",
		Description: "Get real-time quotes for one or more ticker symbols from Tiger Brokers.",
	}, getQuoteHandler(fetcher))
}

func getQuoteHandler(fetcher tiger.QuoteFetcher) mcp.ToolHandlerFor[GetQuoteInput, GetQuoteOutput] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input GetQuoteInput) (*mcp.CallToolResult, GetQuoteOutput, error) {
		if len(input.Symbols) == 0 {
			return nil, GetQuoteOutput{}, fmt.Errorf("get_quote: symbols must contain at least one ticker")
		}

		quotes, err := fetcher.RealTimeQuotes(ctx, input.Symbols)
		if err != nil {
			return nil, GetQuoteOutput{}, fmt.Errorf("get_quote: %w", err)
		}

		return nil, GetQuoteOutput{Quotes: quotes}, nil
	}
}
