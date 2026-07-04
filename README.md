# tiger-mcp

An MCP (Model Context Protocol) server exposing **read-only** Tiger Open API
market data and account access to MCP clients (Claude Desktop, Claude Code,
or any other MCP-capable client). See [specs/mission.md](specs/mission.md),
[specs/tech-stack.md](specs/tech-stack.md), and [specs/roadmap.md](specs/roadmap.md)
for the project's constitution and implementation plan.

This server never places, modifies, or cancels orders — it only reads market
data and your account's positions/assets/order history.

## Prerequisites

- Go 1.25 or later (see `go.mod`).
- A Tiger Open Platform account with API access enabled: TigerID, an RSA
  private key, and a trading account number. See
  [Tiger Open Platform](https://quant.itigerup.com/openapi/) for how to
  register and generate these.

## Build

```sh
git clone <this-repo-url>
cd tiger-mcp
go build -o tiger-mcp ./cmd/tiger-mcp
```

## Configuration

The server reads its Tiger credentials from environment variables at
startup. None of these are ever persisted to disk by the server.

| Variable            | Required | Description                                                        |
| ------------------- | -------- | -------------------------------------------------------------------|
| `TIGER_ID`           | yes      | Your Tiger Open Platform developer/tiger ID.                      |
| `TIGER_PRIVATE_KEY`  | yes      | Your RSA private key, as a raw PEM string (not a file path).       |
| `TIGER_ACCOUNT`      | yes      | Your Tiger trading account number.                                 |

If any of these are missing or empty, the server logs a clear error naming
the specific missing variable and exits immediately — it will not start with
incomplete credentials.

## Run

```sh
TIGER_ID="your-tiger-id" \
TIGER_PRIVATE_KEY="$(cat /path/to/your/private_key.pem)" \
TIGER_ACCOUNT="your-account-number" \
./tiger-mcp
```

The server communicates over stdio using the MCP protocol; it's meant to be
launched as a subprocess by an MCP client, not run interactively.

## Using it with Claude Desktop / Claude Code

Add an entry like this to your MCP client's configuration (e.g. Claude
Desktop's `claude_desktop_config.json`), pointing `command` at your built
binary:

```json
{
  "mcpServers": {
    "tiger-mcp": {
      "command": "/absolute/path/to/tiger-mcp",
      "env": {
        "TIGER_ID": "your-tiger-id",
        "TIGER_PRIVATE_KEY": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
        "TIGER_ACCOUNT": "your-account-number"
      }
    }
  }
}
```

## Available tools

All tools are read-only.

| Tool | Description |
| ---- | ----------- |
| `get_quote` | Get real-time quotes for one or more ticker symbols from Tiger Brokers. |
| `get_market_status` | Get the current trading status (open/closed/etc.) for a Tiger Brokers market. |
| `get_kline` | Get historical K-line (bar) data for a single ticker symbol from Tiger Brokers. |
| `get_depth` | Get order book depth (bids/asks) for one or more ticker symbols from Tiger Brokers. |
| `get_account_assets` | Get the operator's Tiger Brokers account asset summary (buying power, cash, net liquidation, P&L). |
| `get_positions` | Get the operator's currently held positions on Tiger Brokers. |
| `get_orders` | Get the operator's order history on Tiger Brokers (read-only; no filters in this version). |

## Development

```sh
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

## Input from stakeholder
- I want to build tigertrade MCP with golang
