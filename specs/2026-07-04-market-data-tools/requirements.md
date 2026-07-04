# Requirements — Remaining market-data tools (Phase 2)

## Scope
Round out read-only market data with three more MCP tools, all in this one
phase/branch: `get_market_status`, `get_kline`, `get_depth`.

This is the roadmap's [Phase 2](../roadmap.md#phase-2--remaining-market-data-tools).

## In scope
- `get_market_status(market)` → wraps `quote.QuoteClient.GetMarketState(market string) ([]model.MarketState, error)`.
- `get_kline(symbol, period)` → wraps `quote.QuoteClient.GetKline(symbol, period string) ([]model.Kline, error)`.
  Note: unlike `get_quote`, this Tiger endpoint takes a **single** symbol, not a batch.
- `get_depth(symbols, market)` → wraps `quote.QuoteClient.GetQuoteDepth(req model.DepthQuoteRequest) ([]model.Depth, error)`.
  This endpoint **is** batched by symbol, plus an optional market filter.
- Extend `internal/tiger.QuoteFetcher` (or add sibling interfaces) so each new
  tool is unit-testable against a fake, same pattern as `get_quote` in Phase 1.
- Input validation for `market` and `period` against the Tiger SDK's own
  exported enums (see Decisions) — reject unknown values before calling Tiger.

## Out of scope
- Account/positions/orders tools — Phase 3.
- Trading — permanently out of scope per [mission.md](../mission.md#explicit-non-goals-for-now).
- Live validation against the real Tiger API — no credentials available yet,
  same as Phase 1. Tracked as a fast-follow, not a merge blocker.
- Any change to `get_quote` itself (Phase 1, already merged).

## Verified Tiger SDK details (via `go doc` / source read against the pulled
module, not assumed)
- `github.com/tigerfintech/openapi-go-sdk/model` exports typed enums we should
  reuse rather than re-inventing our own:
  - `type Market string` with `MarketAll="ALL"`, `MarketUS="US"`, `MarketHK="HK"`,
    `MarketCN="CN"`, `MarketSG="SG"`.
  - `type BarPeriod string` with `BarPeriodDay="day"`, `BarPeriodWeek="week"`,
    `BarPeriodMonth="month"`, `BarPeriodYear="year"`, `BarPeriod1Min="1min"`,
    `BarPeriod5Min="5min"`, `BarPeriod15Min="15min"`, `BarPeriod30Min="30min"`,
    `BarPeriod60Min="60min"`.
- `model.MarketState{Market, MarketStatus, Status, OpenTime}`.
- `model.Kline{Symbol, Period, NextPageToken, Items []KlineItem}`,
  `model.KlineItem{Time, Volume, Open, Close, High, Low, Amount}`.
- `model.DepthQuoteRequest{Symbols []string, Market, TradeSession, Lang}`.
- `model.Depth{Symbol, Asks []DepthLevel, Bids []DepthLevel}`,
  `model.DepthLevel{Price, Count, Volume}`.

## Decisions
- **Bundle all three tools** in one phase/branch — they're all thin, similarly
  shaped read-only wrappers; low risk to ship together (per roadmap.md's
  Phase 2 grouping).
- **Validate market/period against the Tiger SDK's own enums**, not a
  hand-maintained list of our own. Concretely: accept the input as a plain
  string, check it against the known `model.Market` / `model.BarPeriod`
  constant values, and reject with a clear error listing valid options if it
  doesn't match. This gives strict validation without a second enum to drift
  out of sync — if Tiger's SDK adds a value, updating our check is a one-line,
  visible diff against the SDK's own constants.
- **`get_kline` takes one symbol, `get_depth` takes many** — tool shapes
  mirror the underlying SDK call shapes exactly (same principle as Phase 1's
  `get_quote`), rather than forcing artificial consistency across tools.
- **Merge bar unchanged from Phase 1:** unit tests (with faked Tiger client)
  + code review are sufficient to merge; live Inspector verification against
  the real Tiger API is a tracked, outstanding follow-up (joining the one
  already logged from Phase 1).

## Context
- Constitution docs: [mission.md](../mission.md), [tech-stack.md](../tech-stack.md),
  [roadmap.md](../roadmap.md).
- Builds on Phase 1 (branch `phase-1-get-quote`, merged to `main`): config
  loading, Tiger quote client, `get_quote` tool.
- Working branch: `phase-2-market-data-tools`.
