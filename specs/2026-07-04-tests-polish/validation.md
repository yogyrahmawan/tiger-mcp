# Validation — Tests & polish (Phase 5)

How to know this phase is done and safe to merge.

## Build & static checks
- [x] `go build ./...` succeeds.
- [x] `go vet ./...` reports no issues.
- [x] `gofmt -l .` reports no files.
- [x] `go test ./...` passes (all existing + new tests, no regressions).

## Testability refactor
- [x] `internal/tiger.Client` holds `quoteAPI`/`tradeAPI` interface types,
      not concrete `*quote.QuoteClient`/`*trade.TradeClient`.
- [x] `NewClient`'s behavior is unchanged — it still constructs the real Tiger
      SDK clients, which now satisfy the new interfaces structurally.
- [x] The interfaces are unexported (`quoteAPI`, `tradeAPI`) — internal test
      seam only, not a public API surface.

## Tiger wrapper unit tests — all 7 methods covered
- [x] `RealTimeQuotes`: happy-path mapping + upstream-error passthrough, against a faked `quoteAPI`.
- [x] `MarketStatus`: happy-path mapping + upstream-error passthrough.
- [x] `Kline`: happy-path mapping + upstream-error passthrough.
- [x] `Depth`: happy-path mapping + upstream-error passthrough.
- [x] `Assets`: happy-path mapping + upstream-error passthrough, against a faked `tradeAPI`.
- [x] `Positions`: happy-path mapping + upstream-error passthrough + empty-result-is-not-an-error.
- [x] `Orders`: happy-path mapping + upstream-error passthrough + empty-result-is-not-an-error.

## CI workflow
- [x] `.github/workflows/ci.yml` exists and runs `go build ./...`,
      `go vet ./...`, a `gofmt -l .` check that fails the job on any output,
      and `go test ./...`.
- [x] The workflow's Go version matches `go.mod`.

## Regression check
- [x] All pre-existing tests in `internal/config` and `internal/tools` still
      pass unchanged.
- [x] MCP Inspector `tools/list` against a built binary (dummy-but-well-formed
      credentials) still shows all 7 tools — confirms the `Client` refactor
      didn't change tool wiring or behavior.

## Deferred / follow-up (tracked, not a merge blocker)
- [ ] **Live check (fast-follow, once credentials are available):** the
      running live-check backlog from Phases 1-4 (`get_quote`,
      `get_market_status`, `get_kline`, `get_depth`, `get_account_assets`,
      `get_positions`, `get_orders` against the real Tiger account) remains
      outstanding. This phase does not add to or resolve it.
      **Status: outstanding** — no live credentials available as of this
      phase's implementation (2026-07-04).

## Definition of done (for this phase's merge)
Mergeable when every checkbox under Build/Testability refactor/Tiger wrapper
tests/CI/Regression check is checked, `go test ./...` passes across the whole
repo, and the pre-existing live-check follow-up is acknowledged as still
outstanding rather than silently dropped.
