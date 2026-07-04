# tiger-mcp

[![CI](https://github.com/yogyrahmawan/tiger-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/yogyrahmawan/tiger-mcp/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/yogyrahmawan/tiger-mcp.svg)](https://pkg.go.dev/github.com/yogyrahmawan/tiger-mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25%2B-00ADD8?logo=go)](go.mod)

An MCP ([Model Context Protocol](https://modelcontextprotocol.io)) server,
written in Go, that lets an MCP client (Claude Desktop, Claude Code, or any
other MCP-capable client) query your own [Tiger Brokers](https://www.tigerbrokers.com/)
account and market data through natural language — **read-only**.

This server never places, modifies, or cancels orders. It only reads market
data and your account's positions, assets, and order history.

> [!IMPORTANT]
> This is an independent, community-built project and is **not affiliated
> with, endorsed by, or supported by Tiger Brokers**. It is provided "as is,"
> with no warranty (see [LICENSE](LICENSE)). Market data and account
> information returned by this tool are for informational purposes only and
> are not financial advice. You are solely responsible for verifying any
> data before acting on it, and for keeping your Tiger API credentials
> secure. Use at your own risk.

## Who is this for

Anyone with a Tiger Brokers account and API access who wants an LLM client to
answer questions like "what's my portfolio worth?" or "get me a quote for
AAPL." Each person runs their **own instance** locally against their **own**
Tiger credentials — this is not a shared or hosted service.

## Prerequisites

- Go 1.25 or later (see [go.mod](go.mod)).
- A Tiger Open Platform account with API access enabled: a TigerID, an RSA
  private key, and a trading account number. See
  [Tiger Open Platform](https://quant.itigerup.com/openapi/) for how to
  register and generate these.

## Install

Clone and build from source:

```sh
git clone https://github.com/yogyrahmawan/tiger-mcp.git
cd tiger-mcp
go build -o tiger-mcp ./cmd/tiger-mcp
```

Or install directly with Go:

```sh
go install github.com/yogyrahmawan/tiger-mcp/cmd/tiger-mcp@latest
```

## Configuration

The server reads its Tiger credentials from environment variables at
startup. **None of these are ever persisted to disk by the server.**

| Variable | Required | Description |
| --- | --- | --- |
| `TIGER_ID` | yes | Your Tiger Open Platform developer/tiger ID. |
| `TIGER_PRIVATE_KEY` | yes | Your RSA private key, as a raw PEM string (not a file path). |
| `TIGER_ACCOUNT` | yes | Your Tiger trading account number. |

If any of these are missing or empty, the server logs a clear error naming
the specific missing variable and exits immediately — it will not start with
incomplete credentials.

**Keep your `TIGER_PRIVATE_KEY` secret.** Treat it like a password: don't
commit it, don't paste it into shared chats or issues, and don't log it.

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
| --- | --- |
| `get_quote` | Get real-time quotes for one or more ticker symbols from Tiger Brokers. |
| `get_market_status` | Get the current trading status (open/closed/etc.) for a Tiger Brokers market. |
| `get_kline` | Get historical K-line (bar) data for a single ticker symbol from Tiger Brokers. |
| `get_depth` | Get order book depth (bids/asks) for one or more ticker symbols from Tiger Brokers. |
| `get_account_assets` | Get the operator's Tiger Brokers account asset summary (buying power, cash, net liquidation, P&L). |
| `get_positions` | Get the operator's currently held positions on Tiger Brokers. |
| `get_orders` | Get the operator's order history on Tiger Brokers (read-only; no filters in this version). |

## Project docs

See [specs/mission.md](specs/mission.md), [specs/tech-stack.md](specs/tech-stack.md),
and [specs/roadmap.md](specs/roadmap.md) for the project's constitution and
implementation history.

## Development

```sh
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

## Contributing

Issues and pull requests are welcome. Since this project is strictly
read-only by design (see [specs/mission.md](specs/mission.md)), contributions
that add trading/order-placement capability will not be accepted.

## License

[MIT](LICENSE)
