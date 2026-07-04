# Requirements — Hardening & operator UX (Phase 4)

## Scope
Make the server pleasant and safe to run daily, per roadmap.md's Phase 4:
consistent error handling/input validation, sensible timeouts/retry,
structured logging, and a real operator-facing README.

This is the roadmap's [Phase 4](../roadmap.md#phase-4--hardening--operator-ux).

## In scope
- **Error handling & input validation pass**: review all 7 existing tools
  (`get_quote`, `get_market_status`, `get_kline`, `get_depth`,
  `get_account_assets`, `get_positions`, `get_orders`) for consistent error
  wrapping/messages and confirm no tool can panic on malformed input.
- **Timeouts & retry**: confirmed (see Verified details below) that the Tiger
  SDK already applies a sensible default timeout and retry policy — this
  phase documents/asserts that choice rather than adding new code.
- **Structured logging**: replace the current plain `*log.Logger` in
  `cmd/tiger-mcp/main.go` with `log/slog`, using real levels (Info/Error) and
  structured key-value fields, writing to stderr (stdout stays reserved for
  the MCP protocol, unchanged from Phase 0).
- **README**: setup instructions, prerequisites, env vars table, a Claude
  client config snippet (stdio launch command), and a list of available tools.
- **Fresh-clone verification**: actually `git clone` this repo into a temp
  directory and follow only the README's steps to prove a newcomer can get a
  running server (build, set dummy-but-well-formed env vars, confirm via MCP
  Inspector) with zero prior context.

## Out of scope
- Any new tool or Tiger endpoint — this phase only hardens what exists
  (Phases 1-3).
- A custom retry/backoff layer — superseded by the "use the SDK's own knobs"
  decision below.
- CI setup, linting beyond `go vet`/`gofmt` — that's Phase 5.
- Trading — permanently out of scope per [mission.md](../mission.md#explicit-non-goals-for-now).

## Verified Tiger SDK details (via source read against the pulled module)
- `config.NewClientConfig` defaults `Timeout` to **15 seconds** if unset
  (`config/client_config.go`'s `defaultTimeout = 15 * time.Second`), applied
  whenever `cfg.Timeout == 0`.
- `client.NewHttpClient` (and `NewQuoteHttpClient`) unconditionally sets
  `retryPolicy: DefaultRetryPolicy()` internally — **5 max retries, 60s max
  total retry time, 1s base backoff, 16s max backoff** — applied automatically
  to every request/retryable API call. There is no opt-out or override plumbed
  through `ClientConfig`'s public `Option`s; it's baked into `NewHttpClient`.
- Net effect: we already get a sensible default timeout and light retry/backoff
  for free, simply by using the SDK as we already do in `internal/tiger.NewClient`.
  No new timeout/retry code is needed to satisfy this roadmap bullet.

## Decisions
- **Logging: `log/slog`**, not a hand-rolled prefix scheme. It's stdlib
  (available since Go 1.21; this module targets 1.25), gives real levels, and
  supports structured fields — no new dependency.
- **Timeouts/retry: rely on the Tiger SDK's own defaults**, not a custom
  wrapper. The SDK's 15s timeout / 5-retry backoff policy is reasonable for a
  single-operator local tool; building a parallel retry layer would duplicate
  behavior the SDK already provides and add maintenance surface for no
  benefit. If the SDK's defaults ever prove wrong in practice, revisit then.
- **Verify the fresh-clone claim for real**: clone to a temp dir and actually
  run the documented steps, rather than just proofreading the README. This is
  the strongest and cheapest way to catch a missing step or wrong command.

## Context
- Constitution docs: [mission.md](../mission.md), [tech-stack.md](../tech-stack.md),
  [roadmap.md](../roadmap.md).
- Builds on Phases 0-3 (all merged to `main`): full read-only tool set
  (market data + account/positions) already implemented and unit-tested.
- Working branch: `phase-4-hardening-operator-ux`.
- The live-credential checks logged as outstanding in Phases 1-3's
  validation.md files remain outstanding — this phase doesn't require live
  credentials either (the fresh-clone check uses dummy-but-well-formed values,
  same as prior phases' MCP Inspector checks).
