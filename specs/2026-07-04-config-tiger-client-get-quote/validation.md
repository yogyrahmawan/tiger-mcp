# Validation — Config + Tiger client + `get_quote` tool (Phase 1)

How to know this phase is done and safe to merge. Per the agreed merge bar,
**live Tiger API access is not required to merge** — unit tests + review are
sufficient. Live verification is a tracked fast-follow.

## Build & static checks
- [x] `go build ./...` succeeds.
- [x] `go vet ./...` reports no issues.
- [x] `gofmt -l .` reports no files.

## Config (`internal/config`)
- [x] Loading with all three env vars set succeeds and returns the correct values.
- [x] Loading with `TIGER_ID` missing fails with an error naming `TIGER_ID` specifically.
- [x] Loading with `TIGER_PRIVATE_KEY` missing fails with an error naming it specifically.
- [x] Loading with `TIGER_ACCOUNT` missing fails with an error naming it specifically.
- [x] No credential value is ever logged (check stderr output in tests/manual run).

## Tiger client wrapper (`internal/tiger`)
- [x] `QuoteFetcher` interface exists and the real implementation is
      constructed from `config.Config` using the **verified** (via `go doc`,
      not assumed) Tiger SDK API.
- [x] Code review confirms the SDK call shape matches what group 1 of
      plan.md actually found in the pulled module source.

## `get_quote` tool (`internal/tools`)
- [x] Unit test: multiple symbols → fake `QuoteFetcher` called with the same
      symbols → tool returns a well-formed result reflecting the fake's data.
- [x] Unit test: empty `symbols` input → tool returns a validation error, no
      call made to the fetcher.
- [x] Unit test: fake `QuoteFetcher` returns an error → tool surfaces it as an
      MCP tool error (not a panic or process crash).

## Wire-up (`cmd/tiger-mcp`)
- [x] With required env vars unset, running the binary fails fast with a
      clear stderr message and non-zero exit — it does not hang or start the
      MCP server loop.
- [x] With env vars set to dummy-but-well-formed values (e.g. any string for
      `TIGER_ID`/`TIGER_ACCOUNT`, a syntactically valid-looking PEM for
      `TIGER_PRIVATE_KEY`), the server starts, and MCP Inspector's
      `tools/list` shows exactly one tool: `get_quote`, with the expected
      input schema.

## Deferred / follow-up (tracked, not a merge blocker)
- [ ] **Live check (fast-follow, once credentials are available):** run
      MCP Inspector's `tools/call get_quote` with real symbols against the
      live Tiger account and confirm real quote data comes back. Record the
      outcome here or in a follow-up issue when performed.
      **Status: outstanding** — not yet run, no live credentials available
      as of this phase's implementation (2026-07-04).

## Definition of done (for this phase's merge)
Mergeable when every checkbox under Build/Config/Tiger client wrapper/`get_quote`
tool/Wire-up is checked, `go test ./...` passes, and the live-check follow-up
is explicitly logged as outstanding (not silently dropped).
