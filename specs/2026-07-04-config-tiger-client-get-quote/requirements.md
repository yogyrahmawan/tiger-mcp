# Requirements — Config + Tiger client + `get_quote` tool (Phase 1)

## Scope
The first real vertical slice: load Tiger credentials from the environment,
construct a Tiger Open API client, and expose exactly **one** MCP tool,
`get_quote`, that returns real-time quotes for one or more symbols.

This is the roadmap's [Phase 1](../roadmap.md#phase-1--config--tiger-client--one-get_quote-tool--first-real-slice).

## In scope
- `internal/config`: load and validate three env vars:
  - `TIGER_ID`
  - `TIGER_PRIVATE_KEY` — **raw PEM string**, supplied directly as the env
    var's value (not a file path).
  - `TIGER_ACCOUNT`
  - Missing or empty values fail fast at startup with a clear, specific error
    (naming which var is missing) — never a generic/opaque failure.
- `internal/tiger`: a thin wrapper around the Tiger Go SDK's quote client,
  hidden behind a small interface (see Decisions) so it can be faked in tests.
- `internal/tools`: one MCP tool, `get_quote`, wired into `cmd/tiger-mcp/main.go`.
- Unit tests for config parsing/validation and for the tools layer (using a
  faked tiger client — no live network calls in tests).

## Out of scope (deferred to later phases)
- Any other market-data tool (`get_market_status`, `get_kline`, `get_depth`) — Phase 2.
- Account/positions/orders tools — Phase 3.
- Any trading (place/cancel order) capability — permanently out of scope per
  [mission.md](../mission.md#explicit-non-goals-for-now).
- Live validation against the real Tiger API — no credentials available yet
  (see Decisions). Tracked as a fast-follow, not a blocker for this phase.
- Retry/backoff, rate limiting, structured logging polish — Phase 4.

## `get_quote` tool contract
- **Input:** `symbols: []string` — one or more tickers in a single market
  (matches the underlying SDK's batch quote call shape).
- **Output:** real-time quote data per symbol, as returned by the Tiger SDK,
  surfaced to the MCP client in a structured tool result.
- **Errors:** invalid/empty `symbols`, and upstream Tiger API errors, must
  surface as clear MCP tool errors — not a crashed process.

## Decisions
- **Private key delivery:** raw PEM in `TIGER_PRIVATE_KEY`. Simplest for a
  single stdio/Claude-config user; no filesystem path to manage. (Supersedes
  the "TBD" note in [tech-stack.md](../tech-stack.md#configuration--secrets).)
- **Tool shape:** `get_quote(symbols []string)`, batched, single market —
  mirrors the SDK's `QuoteRealTime(symbols []string)` method directly rather
  than forcing one call per ticker.
- **Testability over live proof, for now:** the operator does not have live
  Tiger credentials available yet. To avoid blocking this phase on that, the
  Tiger client is wrapped behind a small interface (e.g. `QuoteFetcher`) so
  `internal/tools` can be unit-tested against a fake. The real SDK-backed
  implementation is still written and wired up — it's just not exercised
  end-to-end against the live API in this phase.
- **Merge bar:** unit tests (config + tools, with the tiger client mocked)
  passing, plus code review, are sufficient to merge. A live end-to-end
  Inspector run happens as a fast-follow once credentials are available, and
  should be logged as a follow-up task rather than silently skipped.

## Context
- Constitution docs: [mission.md](../mission.md), [tech-stack.md](../tech-stack.md),
  [roadmap.md](../roadmap.md).
- Builds on Phase 0 (branch `phase-0-repo-scaffold`, now merged to `main`):
  empty stdio MCP server, verified via MCP Inspector.
- Working branch: `phase-1-get-quote`.
- **Open risk:** exact Tiger Go SDK method/type names used in earlier
  research (`config.NewClientConfig`, `quote.NewQuoteClient`, `QuoteRealTime`,
  etc.) came from web fetches and are **not yet verified against the actual
  pulled module source**. The plan includes an explicit step to `go get` the
  module and confirm real signatures via `go doc` before writing code against
  them — the same approach used successfully in Phase 0 for the MCP SDK.
