# Validation — Account & positions tools (Phase 3)

How to know this phase is done and safe to merge. Per the agreed merge bar
(same as Phases 1-2), **live Tiger API access is not required to merge** —
unit tests + review are sufficient. Live verification is a tracked fast-follow.

## Build & static checks
- [x] `go build ./...` succeeds.
- [x] `go vet ./...` reports no issues.
- [x] `gofmt -l .` reports no files.

## Safety check — no trading capability introduced
- [x] `grep -rn "PlaceOrder\|ModifyOrder\|CancelOrder\|PlaceForexOrder\|OptionExerciseSubmit\|TransferPosition\|TransferSegmentFund" internal/ cmd/`
      returns **no matches**. This phase adds a `trade.TradeClient`, and it is
      critical that only `Assets`, `Positions`, and `Orders` are ever called
      on it — no mutating/trading method may be wired up, referenced, or
      stubbed anywhere in this codebase.

## Tiger client wrapper (`internal/tiger`)
- [x] `AssetsFetcher`, `PositionsFetcher`, `OrdersFetcher` (or equivalently
      named) interfaces exist, backed by the real Tiger SDK calls confirmed
      in requirements.md (`Assets`, `Positions`, `Orders`), using
      `trade.NewTradeClientFromConfig(client.NewHttpClient(cfg), cfg)`.

## `get_account_assets` tool
- [x] Unit test: happy path → fake fetcher's assets returned as the tool's
      mapped output.
- [x] Unit test: fake fetcher returns an error → surfaced as a tool error.

## `get_positions` tool
- [x] Unit test: happy path with one or more positions → mapped correctly.
- [x] Unit test: empty positions list is treated as a valid, non-error result.
- [x] Unit test: fake fetcher returns an error → surfaced as a tool error.

## `get_orders` tool
- [x] Unit test: happy path with one or more orders → mapped correctly.
- [x] Unit test: empty orders list is treated as a valid, non-error result.
- [x] Unit test: fake fetcher returns an error → surfaced as a tool error.

## Wire-up (`cmd/tiger-mcp`)
- [x] With dummy-but-well-formed credentials, MCP Inspector's `tools/list`
      shows exactly **seven** tools: `get_quote`, `get_market_status`,
      `get_kline`, `get_depth`, `get_account_assets`, `get_positions`,
      `get_orders`, each with the expected input schema.
- [x] Fail-fast-on-missing-env-var behavior (established in Phase 1) still
      holds — no regression from adding these tools.

## Deferred / follow-up (tracked, not a merge blocker)
- [ ] **Live check (fast-follow, once credentials are available):** run MCP
      Inspector `tools/call` for `get_account_assets`, `get_positions`, and
      `get_orders` against the live Tiger account and confirm real data comes
      back, reconciling against the Tiger app per mission.md's success
      criteria. This joins the still-outstanding live checks from Phases 1-2
      (`get_quote`, `get_market_status`, `get_kline`, `get_depth`) — do all
      seven together once credentials exist.
      **Status: outstanding** — not yet run, no live credentials available
      as of this phase's implementation (2026-07-04).

## Definition of done (for this phase's merge)
Mergeable when every checkbox under Build/Safety check/Tiger client
wrapper/all three tools/Wire-up is checked, `go test ./...` passes, and the
live-check follow-up is explicitly logged as outstanding (not silently
dropped).
