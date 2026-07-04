# Plan — Hardening & operator UX (Phase 4)

Coarse, numbered task groups. Each should leave the repo building and (from
group 1 onward) with passing tests.

## 1. Error handling & input validation review
- Re-read all 7 tool handlers (`internal/tools/get_*.go`) for:
  - Consistent error message format (`"<tool_name>: ..."`, matching what's
    already used).
  - No possible panic on malformed/empty input (nil slices, empty strings,
    etc.) — confirm existing validation (`get_quote`, `get_kline`, `get_depth`,
    `get_market_status`) already covers this, and the no-input tools
    (`get_account_assets`, `get_positions`, `get_orders`) have nothing to
    validate.
  - Fix any inconsistency found; add a unit test if a gap is found (unlikely
    given Phases 1-3's coverage, but confirm rather than assume).

## 2. Timeouts & retry — document, don't duplicate
- Add a short code comment in `internal/tiger.NewClient` noting that the
  Tiger SDK's `client.NewHttpClient`/`NewQuoteHttpClient` already apply a 15s
  default timeout and `DefaultRetryPolicy()` (5 retries, 60s max, 1-16s
  backoff) automatically — so this is intentionally not re-implemented here.
- No new retry/timeout code. This group is a documentation/comment change only.

## 3. Structured logging with log/slog
- In `cmd/tiger-mcp/main.go`, replace the `log.New(os.Stderr, "tiger-mcp: ",
  log.LstdFlags)` logger with an `slog.Logger` backed by
  `slog.NewTextHandler(os.Stderr, nil)` (or `NewJSONHandler` — decide during
  implementation based on readability for a single local operator; text is
  likely friendlier for a stdio-run CLI tool).
- Replace `logger.Fatalf(...)` call sites with `logger.Error(...)` followed by
  `os.Exit(1)` (slog has no built-in Fatal), preserving today's fail-fast
  behavior and message content, now as structured key-value pairs (e.g.
  `logger.Error("config error", "err", err)`).

## 4. README rewrite for real operator setup
- Prerequisites: Go version (matches `go.mod`'s `go 1.25.0` / toolchain).
- Clone & build steps.
- Env vars table: `TIGER_ID`, `TIGER_PRIVATE_KEY` (raw PEM), `TIGER_ACCOUNT` —
  what each is, where to get it (Tiger Open Platform), and that none are
  persisted by the server (per mission.md).
- Claude client config snippet: a stdio MCP server entry (command + args +
  env) usable in Claude Desktop/Code config, referencing the built binary path.
- List of available tools (all 7) with a one-line description each, pulled
  from each tool's existing `Description` field so it can't drift.
- Keep the existing "Input from stakeholder" section (repo history) — don't
  delete it, just add the operator-facing content around it.

## 5. Fresh-clone verification
- `git clone` the local repo into a scratch temp directory (or use a bare
  clone of the current branch) so the check reflects only what's committed.
- Follow the new README's steps literally: build, export dummy-but-well-formed
  `TIGER_ID`/`TIGER_PRIVATE_KEY`/`TIGER_ACCOUNT`, run the binary, confirm via
  MCP Inspector `tools/list` that all 7 tools appear — with **zero** knowledge
  assumed beyond what the README states.
- Fix any gap the walkthrough surfaces (missing step, wrong path, unclear
  instruction) before considering this phase done.

## 6. Tests, vet, and final check
- `go test ./...`, `go vet ./...`, `gofmt -l .` all clean.
- Confirm the fail-fast-on-missing-env-var and 7-tool MCP Inspector checks
  from Phase 3 still pass unchanged (no regression from the logging swap).
