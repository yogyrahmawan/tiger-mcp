package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

// GetPositionsInput is the input schema for the get_positions tool. It takes
// no parameters.
type GetPositionsInput struct{}

// GetPositionsOutput is the output schema for the get_positions tool.
type GetPositionsOutput struct {
	Positions []tiger.Position `json:"positions"`
}

// RegisterGetPositions adds the get_positions tool to server, backed by fetcher.
func RegisterGetPositions(server *mcp.Server, fetcher tiger.PositionsFetcher) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_positions",
		Description: "Get the operator's currently held positions on Tiger Brokers.",
	}, getPositionsHandler(fetcher))
}

func getPositionsHandler(fetcher tiger.PositionsFetcher) mcp.ToolHandlerFor[GetPositionsInput, GetPositionsOutput] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ GetPositionsInput) (*mcp.CallToolResult, GetPositionsOutput, error) {
		positions, err := fetcher.Positions(ctx)
		if err != nil {
			return nil, GetPositionsOutput{}, fmt.Errorf("get_positions: %w", err)
		}

		return nil, GetPositionsOutput{Positions: positions}, nil
	}
}
