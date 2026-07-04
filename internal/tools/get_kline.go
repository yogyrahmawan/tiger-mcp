package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

// GetKlineInput is the input schema for the get_kline tool.
type GetKlineInput struct {
	Symbol string `json:"symbol" jsonschema:"ticker symbol, e.g. AAPL"`
	Period string `json:"period" jsonschema:"bar period, one of day, week, month, year, 1min, 5min, 15min, 30min, 60min"`
}

// GetKlineOutput is the output schema for the get_kline tool.
type GetKlineOutput struct {
	Kline tiger.Kline `json:"kline"`
}

// RegisterGetKline adds the get_kline tool to server, backed by fetcher.
func RegisterGetKline(server *mcp.Server, fetcher tiger.KlineFetcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_kline",
		Description: "Get historical K-line (bar) data for a single ticker symbol from Tiger Brokers.",
	}, getKlineHandler(fetcher))
}

func getKlineHandler(fetcher tiger.KlineFetcher) mcp.ToolHandlerFor[GetKlineInput, GetKlineOutput] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input GetKlineInput) (*mcp.CallToolResult, GetKlineOutput, error) {
		if input.Symbol == "" {
			return nil, GetKlineOutput{}, fmt.Errorf("get_kline: symbol is required")
		}
		if !tiger.IsValidBarPeriod(input.Period) {
			return nil, GetKlineOutput{}, fmt.Errorf(
				"get_kline: invalid period %q, must be one of: %s",
				input.Period, strings.Join(tiger.ValidBarPeriods(), ", "),
			)
		}

		kline, err := fetcher.Kline(ctx, input.Symbol, input.Period)
		if err != nil {
			return nil, GetKlineOutput{}, fmt.Errorf("get_kline: %w", err)
		}

		return nil, GetKlineOutput{Kline: *kline}, nil
	}
}
