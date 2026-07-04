# Mission

## What we are building
`tiger-mcp` is a **Model Context Protocol (MCP) server**, written in Go, that lets an
LLM client (Claude Desktop, Claude Code, or any MCP-capable client) query a
**Tiger Brokers (Tiger Open API)** account and market data through natural language.

## Who it is for
A **single user querying their own live Tiger account** — a personal trading
assistant. The operator runs the server locally and connects it to their own
Claude client. This is not (yet) a multi-tenant or publicly distributed service,
so we optimize for one user's ergonomics over generality.

## What it does
Read-only access to two capability areas, exposed as MCP tools:

1. **Market data** — real-time quotes, market status, K-line history, depth.
2. **Account & positions** — account assets, open positions, order history.

## Explicit non-goals (for now)
- **No trading.** The server does **not** place, modify, or cancel orders. It is
  strictly read-only. An LLM must never be able to move real money through it.
- **No multi-user / hosting.** Runs as a local stdio subprocess for one operator.
- **No credential storage.** Credentials are supplied by the environment at runtime;
  the server never persists them.

## Success looks like
- The operator can ask their Claude client "what's my Tiger portfolio worth?" or
  "get me a quote for AAPL" and get accurate, live answers.
- Adding a new read-only Tiger endpoint is a small, mechanical change (one tool).
- A wrong or missing credential fails fast with a clear, actionable message.

## Guiding principles
- **Read-only by construction.** Trading endpoints of the Tiger SDK are simply not
  wired up. Safety comes from what we don't build, not from runtime guards alone.
- **Thin, honest tools.** Each MCP tool maps closely to one Tiger endpoint with
  clear inputs/outputs; no hidden trading side effects.
- **Fail loud, fail early.** Validate config and connectivity at startup.
- **Small vertical slices.** Ship the thinnest end-to-end path first, then widen.
