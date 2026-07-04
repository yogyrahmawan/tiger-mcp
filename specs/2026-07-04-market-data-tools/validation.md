# Validation — Remaining market-data tools (Phase 2)

How to know this phase is done and safe to merge. Per the agreed merge bar
(same as Phase 1), **live Tiger API access is not required to merge** — unit
tests + review are sufficient. Live verification is a tracked fast-follow.

## Build & static checks
- [x] `go build ./...` succeeds.
- [x] `go vet ./...` reports no issues.
- [x] `gofmt -l .` reports no files.

## Tiger client wrapper (`internal/tiger`)
- [x] `MarketStatusFetcher`, `KlineFetcher`, `DepthFetcher` (or equivalently
      named) interfaces exist, backed by the real Tiger SDK calls confirmed
      in requirements.md (`GetMarketState`, `GetKline`, `GetQuoteDepth`).
- [x] Market/period validation helpers check against the Tiger SDK's own
      `model.Market*` / `model.BarPeriod*` constants — not a separately
      hand-maintained list that could drift.

## `get_market_status` tool
- [x] Unit test: valid market → fake fetcher called with that market → tool
      returns the fake's mapped result.
- [x] Unit test: invalid/unknown market string → validation error, fetcher
      not called.
- [x] Unit test: fake fetcher returns an error → surfaced as a tool error.

## `get_kline` tool
- [x] Unit test: valid symbol + period → fake fetcher called correctly →
      well-formed result.
- [x] Unit test: invalid period → validation error, fetcher not called.
- [x] Unit test: empty symbol → validation error, fetcher not called.
- [x] Unit test: fake fetcher returns an error → surfaced as a tool error.

## `get_depth` tool
- [x] Unit test: multiple symbols (+ optional market) → fake fetcher called
      correctly → well-formed result.
- [x] Unit test: empty symbols → validation error, fetcher not called.
- [x] Unit test: invalid market (when provided) → validation error, fetcher
      not called.
- [x] Unit test: fake fetcher returns an error → surfaced as a tool error.

## Wire-up (`cmd/tiger-mcp`)
- [x] With dummy-but-well-formed credentials, MCP Inspector's `tools/list`
      shows exactly **four** tools: `get_quote`, `get_market_status`,
      `get_kline`, `get_depth`, each with the expected input schema.
- [x] Fail-fast-on-missing-env-var behavior (established in Phase 1) still
      holds — no regression from adding these tools.

## Deferred / follow-up (tracked, not a merge blocker)
- [ ] **Live check (fast-follow, once credentials are available):** run MCP
      Inspector `tools/call` for `get_market_status`, `get_kline`, and
      `get_depth` with real inputs against the live Tiger account and confirm
      real data comes back. This joins the still-outstanding Phase 1 live
      check for `get_quote` — do them together once credentials exist, rather
      than tracking two separate lists.
      **Status: outstanding** — not yet run, no live credentials available
      as of this phase's implementation (2026-07-04).

## Definition of done (for this phase's merge)
Mergeable when every checkbox under Build/Tiger client wrapper/all three
tools/Wire-up is checked, `go test ./...` passes, and the live-check
follow-up is explicitly logged as outstanding (not silently dropped).
