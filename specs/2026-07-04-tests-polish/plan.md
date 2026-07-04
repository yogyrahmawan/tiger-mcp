# Plan — Tests & polish (Phase 5)

Coarse, numbered task groups. Each should leave the repo building and (from
group 2 onward) with passing tests.

## 1. Introduce quoteAPI/tradeAPI interfaces as a test seam
- In `internal/tiger`, define:
  ```go
  type quoteAPI interface {
      GetBrief(req model.BriefRequest) ([]model.Brief, error)
      GetMarketState(market string) ([]model.MarketState, error)
      GetKline(symbol, period string) ([]model.Kline, error)
      GetQuoteDepth(req model.DepthQuoteRequest) ([]model.Depth, error)
  }

  type tradeAPI interface {
      Assets(req model.AssetsRequest) ([]model.Asset, error)
      Positions(req model.PositionsRequest) ([]model.Position, error)
      Orders(req model.OrdersRequest) ([]model.Order, error)
  }
  ```
- Change `Client`'s fields from `*quote.QuoteClient`/`*trade.TradeClient` to
  `quoteAPI`/`tradeAPI`. `NewClient` is otherwise unchanged — the real SDK
  clients already satisfy these interfaces structurally.
- `go build ./...` must still succeed with zero other changes — this group is
  a pure refactor, no behavior change.

## 2. Unit tests for quote-side Client methods
- New `internal/tiger/tiger_test.go` (or split per file, matching existing
  `tiger.go`/`market_data.go` structure — decide during implementation) with
  fake `quoteAPI` implementations.
- Cover `RealTimeQuotes`, `MarketStatus`, `Kline`, `Depth`: happy-path mapping
  from a fake SDK response to our curated struct, and upstream-error passthrough.

## 3. Unit tests for trade-side Client methods
- Fake `tradeAPI` implementation.
- Cover `Assets`, `Positions`, `Orders`: happy-path mapping and upstream-error
  passthrough, including an empty-result case (valid, not an error) mirroring
  the tool-level tests already in `internal/tools`.

## 4. CI workflow
- Add `.github/workflows/ci.yml`: on push and pull_request to any branch (or
  at least `main`), set up Go (matching `go.mod`'s version), then run
  `go build ./...`, `go vet ./...`, a `gofmt -l .` step that fails if it
  prints anything, and `go test ./...`.

## 5. Final polish pass
- `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...` all clean
  across the whole repo (not just `internal/tiger`).
- Skim `internal/tools`' existing tests once more for consistency with the
  new `internal/tiger` tests (naming, fake-struct patterns) — align style if
  there's drift, but don't rewrite working tests without a reason.
