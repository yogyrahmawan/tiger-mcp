# Validation — Hardening & operator UX (Phase 4)

How to know this phase is done and safe to merge.

## Build & static checks
- [ ] `go build ./...` succeeds.
- [ ] `go vet ./...` reports no issues.
- [ ] `gofmt -l .` reports no files.
- [ ] `go test ./...` passes (all existing tests, no regressions).

## Error handling & input validation
- [ ] All 7 tool handlers reviewed; error messages consistently prefixed
      with `"<tool_name>: ..."`.
- [ ] No tool handler can panic on empty/nil/malformed input (confirmed by
      inspection, backed by existing or newly-added unit tests).

## Timeouts & retry
- [ ] `internal/tiger.NewClient` has a comment documenting that the Tiger
      SDK's default 15s timeout and `DefaultRetryPolicy()` (5 retries, 60s
      max, 1-16s backoff) are relied upon intentionally — no duplicate
      custom retry/timeout logic was added.

## Structured logging
- [ ] `cmd/tiger-mcp/main.go` uses `log/slog` instead of a plain `*log.Logger`.
- [ ] Fatal error paths (`config.Load` failure, `tiger.NewClient` failure,
      `server.Run` failure) log at Error level with structured fields (e.g.
      `"err"`) and still exit non-zero, preserving Phase 0-3's fail-fast
      behavior and message clarity.
- [ ] Nothing is ever written to **stdout** by the logger (stdout remains
      reserved for the MCP protocol stream).

## README
- [ ] Documents prerequisites (Go version), clone/build steps, all three env
      vars with descriptions, a working Claude client stdio config snippet,
      and a list of all 7 tools.
- [ ] No credential values or examples that look like real secrets are
      included.

## Fresh-clone verification (performed for real, not just reviewed)
- [ ] Cloned the repo into a scratch temp directory.
- [ ] Followed only the README's steps: build succeeded as documented.
- [ ] Set dummy-but-well-formed `TIGER_ID`/`TIGER_PRIVATE_KEY`/`TIGER_ACCOUNT`
      per the README's instructions; server started without needing any
      undocumented step.
- [ ] MCP Inspector's `tools/list` against the freshly-built binary shows all
      **7** tools, matching Phase 3's baseline (no regression).
- [ ] Any gap found during this walkthrough was fixed in the README before
      checking this box (not deferred).

## Deferred / follow-up (tracked, not a merge blocker)
- [ ] **Live check (fast-follow, once credentials are available):** the
      running live-check backlog from Phases 1-3 (`get_quote`,
      `get_market_status`, `get_kline`, `get_depth`, `get_account_assets`,
      `get_positions`, `get_orders` against the real Tiger account) remains
      outstanding. This phase does not add to or resolve it — noted here only
      to keep it visible.
      **Status: outstanding** — no live credentials available as of this
      phase's implementation (2026-07-04).

## Definition of done (for this phase's merge)
Mergeable when every checkbox under Build/Error handling/Timeouts/Logging/README
is checked, the fresh-clone verification was actually performed (not just
reviewed) and passed, and the pre-existing live-check follow-up is
acknowledged as still outstanding rather than silently dropped.
