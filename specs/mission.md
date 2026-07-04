# Mission

## What we are building
`tiger-mcp` is a **Model Context Protocol (MCP) server**, written in Go, that lets an
LLM client (Claude Desktop, Claude Code, or any MCP-capable client) query a
**Tiger Brokers (Tiger Open API)** account and market data through natural language.

## Who it is for
Anyone with a Tiger Brokers account who wants their own MCP client (Claude
Desktop, Claude Code, or any other MCP-capable client) to query their **own
live Tiger account** — a personal trading assistant, distributed as an
open-source tool so others can run it against their own credentials. Each
operator runs their own instance locally; there is no shared/hosted service.
It's still fundamentally a **single-operator, single-account** tool per
instance — "public" means the source and setup are open to anyone, not that
one deployment serves multiple users.

## What it does
Read-only access to two capability areas, exposed as MCP tools:

1. **Market data** — real-time quotes, market status, K-line history, depth.
2. **Account & positions** — account assets, open positions, order history.

## Explicit non-goals (for now)
- **No trading.** The server does **not** place, modify, or cancel orders. It is
  strictly read-only. An LLM must never be able to move real money through it.
- **No multi-user / hosting.** Each instance runs as a local stdio subprocess
  for one operator and one Tiger account. Being open-source doesn't change
  this — there's no shared/hosted deployment serving multiple users.
- **No credential storage.** Credentials are supplied by the environment at runtime;
  the server never persists them.

## Success looks like
- The operator can ask their Claude client "what's my Tiger portfolio worth?" or
  "get me a quote for AAPL" and get accurate, live answers.
- Adding a new read-only Tiger endpoint is a small, mechanical change (one tool).
- A wrong or missing credential fails fast with a clear, actionable message.
- A new operator can clone the repo and get a running server using only the
  README — no undocumented steps, no tribal knowledge.

## Guiding principles
- **Read-only by construction.** Trading endpoints of the Tiger SDK are simply not
  wired up. Safety comes from what we don't build, not from runtime guards alone.
- **Thin, honest tools.** Each MCP tool maps closely to one Tiger endpoint with
  clear inputs/outputs; no hidden trading side effects.
- **Fail loud, fail early.** Validate config and connectivity at startup.
- **Small vertical slices.** Ship the thinnest end-to-end path first, then widen.
