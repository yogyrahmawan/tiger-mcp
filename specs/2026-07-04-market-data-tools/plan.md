# Plan — Remaining market-data tools (Phase 2)

Coarse, numbered task groups. Each should leave the repo building and (from
group 2 onward) with passing tests. SDK signatures below are already
verified against the pulled `openapi-go-sdk@v0.4.2` source (see
requirements.md) — no re-verification step needed this time, unlike Phase 1.

## 1. Extend `internal/tiger` with three new fetch methods
- Add to (or alongside) the existing `QuoteFetcher`-style interfaces:
  ```go
  type MarketStatusFetcher interface {
      MarketStatus(ctx context.Context, market string) ([]MarketState, error)
  }
  type KlineFetcher interface {
      Kline(ctx context.Context, symbol, period string) (*Kline, error)
  }
  type DepthFetcher interface {
      Depth(ctx context.Context, symbols []string, market string) ([]Depth, error)
  }
  ```
- Define local `MarketState`, `Kline`/`KlineItem`, `Depth`/`DepthLevel` structs
  (mirroring the Tiger SDK's model types, same pattern as `tiger.Quote` in Phase 1).
- Implement each against `*Client` using the verified SDK calls:
  `GetMarketState(market)`, `GetKline(symbol, period)`, `GetQuoteDepth(model.DepthQuoteRequest{...})`.
- Add a small validated-enum helper, e.g. `isValidMarket(s string) bool` /
  `isValidBarPeriod(s string) bool`, checking against the Tiger SDK's own
  `model.Market...` / `model.BarPeriod...` constants.

## 2. `get_market_status` tool
- Input: `market string`. Validate against known `model.Market` values before calling.
- Output: the mapped `[]MarketState`.
- Unit tests: valid market happy path, invalid market rejected without calling
  the fetcher, upstream error passthrough.

## 3. `get_kline` tool
- Input: `symbol string`, `period string`. Validate `period` against known
  `model.BarPeriod` values; require non-empty `symbol`.
- Output: the mapped `Kline` (symbol + items).
- Unit tests: happy path, invalid period rejected, empty symbol rejected,
  upstream error passthrough.

## 4. `get_depth` tool
- Input: `symbols []string`, `market string` (optional — empty means "let
  Tiger apply its default", matching `DepthQuoteRequest.Market`'s `omitempty`).
- Output: the mapped `[]Depth`.
- Unit tests: happy path (multiple symbols), empty-symbols rejected, upstream
  error passthrough. If `market` is provided, validate it against known
  `model.Market` values same as group 2.

## 5. Wire-up (`cmd/tiger-mcp/main.go`)
- Register all three new tools alongside `get_quote` via
  `tools.RegisterGetMarketStatus`, `tools.RegisterGetKline`, `tools.RegisterGetDepth`
  (or one `tools.RegisterAll(server, tigerClient)` convenience function —
  whichever keeps `main.go` cleanest, decide during implementation).

## 6. Tests, vet, and follow-up logging
- `go test ./...`, `go vet ./...`, `gofmt -l .` all clean.
- Run the same fail-fast / MCP Inspector `tools/list` check as Phase 1
  (dummy-but-well-formed credentials), now expecting **four** tools total:
  `get_quote`, `get_market_status`, `get_kline`, `get_depth`.
- Update validation.md's deferred section to note these three new tools also
  need a live check once credentials are available — don't create a second,
  disconnected follow-up list from Phase 1's.
