# Requirements — Account & positions tools (Phase 3)

## Scope
Portfolio read access: three MCP tools, `get_account_assets`, `get_positions`,
`get_orders`, all shipping together in this one phase/branch (matching the
Phase 2 precedent).

This is the roadmap's [Phase 3](../roadmap.md#phase-3--account--positions-tools).

## In scope
- `get_account_assets()` → wraps `trade.TradeClient.Assets(model.AssetsRequest{}) ([]model.Asset, error)`.
- `get_positions()` → wraps `trade.TradeClient.Positions(model.PositionsRequest{}) ([]model.Position, error)`.
- `get_orders()` → wraps `trade.TradeClient.Orders(model.OrdersRequest{}) ([]model.Order, error)`.
  **No filter params in v1** — a plain call returning Tiger's default order
  history, per the agreed decision below.
- Extend `internal/tiger` with a `TradeClient` (or reuse `Client` — decide
  during implementation) wrapping the Tiger SDK's `trade.TradeClient`,
  constructed alongside the existing quote client.
- Unit tests for all three tools against faked fetchers, same pattern as
  Phases 1-2.

## Explicitly out of scope — permanent, not just deferred
- **No trading.** `trade.TradeClient` also exposes `PlaceOrder`, `ModifyOrder`,
  `CancelOrder`, and several other write/mutating methods (forex orders,
  option exercise, position transfers, segment fund transfers). **None of
  these are to be called anywhere in this codebase.** This is not a
  this-phase omission — it is the mission's core safety guarantee (see
  [mission.md](../mission.md#explicit-non-goals-for-now)). Only `Assets`,
  `Positions`, and `Orders` (all read-only queries) are wired up.
- Order filtering by status/date-range/symbol — `OrdersRequest` supports
  `StartDate`, `EndDate`, `States`, `Limit`, etc., but v1 passes an empty
  request and takes Tiger's defaults. Can be added later as tool input params
  without changing the tool's shape.
- Live validation against the real Tiger API — no credentials available yet,
  same as Phases 1-2. Tracked as a fast-follow, not a merge blocker.

## Verified Tiger SDK details (via source read against the pulled module,
not assumed — see `examples/trade/main.go` for confirmed real usage)
- Trade client construction differs from the quote client:
  `trade.NewTradeClientFromConfig(client.NewHttpClient(cfg), cfg)` —
  note **`client.NewHttpClient`**, not `NewQuoteHttpClient` (that's
  quote-specific). Passing the account via `cfg` means request structs can be
  left empty (`model.AssetsRequest{}` etc.) and Tiger infers the account.
- `model.Asset{Account, Capability, Currency, BuyingPower, CashValue,
  NetLiquidation, RealizedPnL, UnrealizedPnL, Segments []AssetSegment}`.
- `model.Position{Account, Symbol, SecType, Market, Currency, Position,
  PositionQty, AverageCost, MarketValue, RealizedPnl, UnrealizedPnl,
  UnrealizedPnlPercent, LatestPrice, Name, ...}` — a large struct; we surface
  a curated subset (see plan.md), not every field.
- `model.Order{Account, ID, OrderId, Action, OrderType, TotalQuantity,
  LimitPrice, Status, FilledQuantity, AvgFillPrice, Symbol, SecType, Market,
  Currency, Commission, RealizedPnl, OpenTime, UpdateTime, ...}` — same
  curation approach.

## Decisions
- **Bundle all three tools** in one phase/branch, consistent with Phase 2.
- **`get_orders` ships with no filter params in v1** — simplest read, matches
  the read-only mission scope; filtering is a natural, additive follow-up.
- **Curate output fields rather than mirroring every SDK field 1:1** —
  `Position` and `Order` have 20-30+ fields each; we expose the ones useful
  for an LLM answering portfolio questions (holdings, cost basis, P&L, order
  status/fill info), not internal/rarely-used ones. Follows the same "thin,
  honest tool" principle as `Quote` in Phase 1, just with more source fields
  to choose from.
- **Merge bar unchanged from Phases 1-2:** unit tests (faked trade client) +
  code review are sufficient to merge; live Inspector verification against
  the real Tiger API remains a tracked, outstanding follow-up — still no
  live credentials available as of this phase.

## Context
- Constitution docs: [mission.md](../mission.md), [tech-stack.md](../tech-stack.md),
  [roadmap.md](../roadmap.md).
- Builds on Phase 1 (`phase-1-get-quote`) and Phase 2
  (`phase-2-market-data-tools`), both merged to `main`.
- Working branch: `phase-3-account-positions`.
