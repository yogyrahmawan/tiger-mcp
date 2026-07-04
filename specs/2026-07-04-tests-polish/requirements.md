# Requirements — Tests & polish (Phase 5)

## Scope
Close the last gap called out in [tech-stack.md](../tech-stack.md#testing--quality):
unit tests for the `internal/tiger` wrapper (config parsing is already
covered from Phase 1). Also add a basic CI workflow, per roadmap.md's
"small CI if desired."

This is the roadmap's [Phase 5](../roadmap.md#phase-5--tests--polish-as-needed) —
the last phase currently on the roadmap.

## In scope
- **Make `internal/tiger.Client` testable** by introducing two small,
  unexported interfaces capturing exactly the Tiger SDK methods it calls:
  - `quoteAPI`: `GetBrief`, `GetMarketState`, `GetKline`, `GetQuoteDepth`
    (confirmed as the complete set of quote-side calls via
    `grep -n "c\.quoteClient\." internal/tiger/*.go`).
  - `tradeAPI`: `Assets`, `Positions`, `Orders` (complete set of trade-side calls).
  `Client` holds these interface types instead of the concrete
  `*quote.QuoteClient`/`*trade.TradeClient`; the real SDK types already
  satisfy them structurally, so `NewClient` needs no behavior change beyond
  the field type declarations.
- **Unit tests for all 7 `Client` methods** (`RealTimeQuotes`, `MarketStatus`,
  `Kline`, `Depth`, `Assets`, `Positions`, `Orders`), each against a faked
  `quoteAPI`/`tradeAPI`, covering: happy-path mapping (SDK model → our
  curated struct) and upstream-error passthrough.
- **A basic CI workflow** (`.github/workflows/ci.yml`) running
  `go build ./...`, `go vet ./...`, `gofmt -l .` (failing on output), and
  `go test ./...` on push/PR — useful now that `main` will eventually be
  pushed to `origin` and reviewed via PRs.

## Out of scope
- Any new tool, endpoint, or behavior change — this phase is tests/CI only.
- Live validation against the real Tiger API — still no credentials
  available; the running Phases 1-4 live-check backlog is untouched by this
  phase.
- Trading — permanently out of scope per [mission.md](../mission.md#explicit-non-goals-for-now).
- Linting beyond `go vet`/`gofmt` (no golangci-lint or similar) — not
  requested, keep the CI workflow minimal.

## Decisions
- **Testability via narrow SDK interfaces**, not extracted pure-mapping
  functions. This matches tech-stack.md's already-stated plan ("SDK client
  mocked behind a small interface") and gives coverage of the actual
  `Client` methods the tools depend on, not just the mapping logic in isolation.
- **Full coverage of all 7 fetcher methods**, not a representative subset —
  this is the phase explicitly dedicated to closing the test gap, so do it
  completely rather than partially now and coming back later.
- **Add CI now** — `main` already has 11 unpushed commits and will be pushed
  to `origin` at some point; a basic build/vet/test workflow costs little and
  catches regressions on future pushes/PRs.
- **Interfaces stay unexported** (`quoteAPI`, `tradeAPI`) — they exist purely
  as an internal test seam, not a public extension point. No other package
  needs to implement them.

## Context
- Constitution docs: [mission.md](../mission.md), [tech-stack.md](../tech-stack.md),
  [roadmap.md](../roadmap.md).
- Builds on Phases 0-4, all merged to `main`. This is the last phase
  currently defined on the roadmap; only "Future (out of current scope)"
  items remain after this (SSE/HTTP transport, sandbox switch, trading —
  none of which are being started here).
- Working branch: `phase-5-tests-polish`.
- Confirmed call sites (via source grep, not assumed):
  - `internal/tiger/tiger.go:76` → `GetBrief`
  - `internal/tiger/market_data.go:68,88,115` → `GetMarketState`, `GetKline`, `GetQuoteDepth`
  - `internal/tiger/account.go:74,96,123` → `Assets`, `Positions`, `Orders`
