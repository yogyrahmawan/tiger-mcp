package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

// GetMarketStatusInput is the input schema for the get_market_status tool.
type GetMarketStatusInput struct {
	Market string `json:"market" jsonschema:"market code to check, one of ALL, US, HK, CN, SG"`
}

// GetMarketStatusOutput is the output schema for the get_market_status tool.
type GetMarketStatusOutput struct {
	MarketStates []tiger.MarketState `json:"marketStates"`
}

// RegisterGetMarketStatus adds the get_market_status tool to server, backed by fetcher.
func RegisterGetMarketStatus(server *mcp.Server, fetcher tiger.MarketStatusFetcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_market_status",
		Description: "Get the current trading status (open/closed/etc.) for a Tiger Brokers market.",
	}, getMarketStatusHandler(fetcher))
}

func getMarketStatusHandler(fetcher tiger.MarketStatusFetcher) mcp.ToolHandlerFor[GetMarketStatusInput, GetMarketStatusOutput] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input GetMarketStatusInput) (*mcp.CallToolResult, GetMarketStatusOutput, error) {
		if !tiger.IsValidMarket(input.Market) {
			return nil, GetMarketStatusOutput{}, fmt.Errorf(
				"get_market_status: invalid market %q, must be one of: %s",
				input.Market, strings.Join(tiger.ValidMarkets(), ", "),
			)
		}

		states, err := fetcher.MarketStatus(ctx, input.Market)
		if err != nil {
			return nil, GetMarketStatusOutput{}, fmt.Errorf("get_market_status: %w", err)
		}

		return nil, GetMarketStatusOutput{MarketStates: states}, nil
	}
}
