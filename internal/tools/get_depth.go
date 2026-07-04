package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

// GetDepthInput is the input schema for the get_depth tool.
type GetDepthInput struct {
	Symbols []string `json:"symbols" jsonschema:"ticker symbols to fetch order book depth for, e.g. AAPL, TSLA"`
	Market  string   `json:"market,omitempty" jsonschema:"optional market filter, one of ALL, US, HK, CN, SG"`
}

// GetDepthOutput is the output schema for the get_depth tool.
type GetDepthOutput struct {
	Depths []tiger.Depth `json:"depths"`
}

// RegisterGetDepth adds the get_depth tool to server, backed by fetcher.
func RegisterGetDepth(server *mcp.Server, fetcher tiger.DepthFetcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_depth",
		Description: "Get order book depth (bids/asks) for one or more ticker symbols from Tiger Brokers.",
	}, getDepthHandler(fetcher))
}

func getDepthHandler(fetcher tiger.DepthFetcher) mcp.ToolHandlerFor[GetDepthInput, GetDepthOutput] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input GetDepthInput) (*mcp.CallToolResult, GetDepthOutput, error) {
		if len(input.Symbols) == 0 {
			return nil, GetDepthOutput{}, fmt.Errorf("get_depth: symbols must contain at least one ticker")
		}
		if input.Market != "" && !tiger.IsValidMarket(input.Market) {
			return nil, GetDepthOutput{}, fmt.Errorf(
				"get_depth: invalid market %q, must be one of: %s",
				input.Market, strings.Join(tiger.ValidMarkets(), ", "),
			)
		}

		depths, err := fetcher.Depth(ctx, input.Symbols, input.Market)
		if err != nil {
			return nil, GetDepthOutput{}, fmt.Errorf("get_depth: %w", err)
		}

		return nil, GetDepthOutput{Depths: depths}, nil
	}
}
