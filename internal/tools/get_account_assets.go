package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

// GetAccountAssetsInput is the input schema for the get_account_assets tool.
// It takes no parameters.
type GetAccountAssetsInput struct{}

// GetAccountAssetsOutput is the output schema for the get_account_assets tool.
type GetAccountAssetsOutput struct {
	Assets []tiger.Asset `json:"assets"`
}

// RegisterGetAccountAssets adds the get_account_assets tool to server, backed by fetcher.
func RegisterGetAccountAssets(server *mcp.Server, fetcher tiger.AssetsFetcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_account_assets",
		Description: "Get the operator's Tiger Brokers account asset summary (buying power, cash, net liquidation, P&L).",
	}, getAccountAssetsHandler(fetcher))
}

func getAccountAssetsHandler(fetcher tiger.AssetsFetcher) mcp.ToolHandlerFor[GetAccountAssetsInput, GetAccountAssetsOutput] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ GetAccountAssetsInput) (*mcp.CallToolResult, GetAccountAssetsOutput, error) {
		assets, err := fetcher.Assets(ctx)
		if err != nil {
			return nil, GetAccountAssetsOutput{}, fmt.Errorf("get_account_assets: %w", err)
		}

		return nil, GetAccountAssetsOutput{Assets: assets}, nil
	}
}
