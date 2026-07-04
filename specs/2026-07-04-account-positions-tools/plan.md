# Plan — Account & positions tools (Phase 3)

Coarse, numbered task groups. Each should leave the repo building and (from
group 2 onward) with passing tests. SDK signatures are already verified
against the pulled `openapi-go-sdk@v0.4.2` source (see requirements.md) — no
re-verification step needed.

## 1. Add a Tiger trade client alongside the existing quote client
- In `internal/tiger`, construct a `trade.TradeClient` at `NewClient` time:
  `trade.NewTradeClientFromConfig(client.NewHttpClient(tigerCfg), tigerCfg)`.
  Store it on the existing `Client` struct alongside `quoteClient` — one
  `Client` backing all fetcher interfaces (quote + trade), consistent with
  how `main.go` already wires a single `tigerClient` into every `Register*` call.
- Define curated local structs: `Asset`, `Position`, `Order` (subset of
  fields per requirements.md's "Decisions" — pick fields useful for an LLM
  answering portfolio questions; do not mirror every SDK field).
- Define fetcher interfaces:
  ```go
  type AssetsFetcher interface {
      Assets(ctx context.Context) ([]Asset, error)
  }
  type PositionsFetcher interface {
      Positions(ctx context.Context) ([]Position, error)
  }
  type OrdersFetcher interface {
      Orders(ctx context.Context) ([]Order, error)
  }
  ```
- Implement each on `*Client` using `tc.Assets(model.AssetsRequest{})`,
  `tc.Positions(model.PositionsRequest{})`, `tc.Orders(model.OrdersRequest{})`.
- **Do not** add any wrapper for `PlaceOrder`/`ModifyOrder`/`CancelOrder` or
  any other mutating trade method — not now, not as a "just in case" stub.

## 2. `get_account_assets` tool
- No input (empty struct `struct{}` or equivalent minimal input type).
- Output: mapped `[]Asset`.
- Unit tests: happy path (fake returns sample assets), upstream error passthrough.

## 3. `get_positions` tool
- No input.
- Output: mapped `[]Position`.
- Unit tests: happy path (including an empty-positions case — a valid "no
  holdings" result, not an error), upstream error passthrough.

## 4. `get_orders` tool
- No input (v1 — no filters, per requirements.md decision).
- Output: mapped `[]Order`.
- Unit tests: happy path (including empty-orders case), upstream error passthrough.

## 5. Wire-up (`cmd/tiger-mcp/main.go`)
- Register all three new tools via `tools.RegisterGetAccountAssets`,
  `tools.RegisterGetPositions`, `tools.RegisterGetOrders`, alongside the four
  existing tools.

## 6. Tests, vet, and follow-up logging
- `go test ./...`, `go vet ./...`, `gofmt -l .` all clean.
- Run the same fail-fast / MCP Inspector `tools/list` check as prior phases
  (dummy-but-well-formed credentials), now expecting **seven** tools total:
  `get_quote`, `get_market_status`, `get_kline`, `get_depth`,
  `get_account_assets`, `get_positions`, `get_orders`.
- Update validation.md's deferred section to fold these three tools into the
  same running live-check follow-up list from Phases 1-2, rather than a new
  disconnected one.
