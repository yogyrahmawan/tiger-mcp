package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

// GetOrdersInput is the input schema for the get_orders tool. It takes no
// parameters in v1 (no filtering by status/date-range/symbol yet).
type GetOrdersInput struct{}

// GetOrdersOutput is the output schema for the get_orders tool.
type GetOrdersOutput struct {
	Orders []tiger.Order `json:"orders"`
}

// RegisterGetOrders adds the get_orders tool to server, backed by fetcher.
func RegisterGetOrders(server *mcp.Server, fetcher tiger.OrdersFetcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_orders",
		Description: "Get the operator's order history on Tiger Brokers (read-only; no filters in this version).",
	}, getOrdersHandler(fetcher))
}

func getOrdersHandler(fetcher tiger.OrdersFetcher) mcp.ToolHandlerFor[GetOrdersInput, GetOrdersOutput] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ GetOrdersInput) (*mcp.CallToolResult, GetOrdersOutput, error) {
		orders, err := fetcher.Orders(ctx)
		if err != nil {
			return nil, GetOrdersOutput{}, fmt.Errorf("get_orders: %w", err)
		}

		return nil, GetOrdersOutput{Orders: orders}, nil
	}
}
