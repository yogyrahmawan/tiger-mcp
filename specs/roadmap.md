# Roadmap

High-level implementation order, in **very small, independently verifiable phases**.
Each phase should be a single focused change set that leaves the repo building.

---

## Phase 0 — Repo scaffold & empty server
**Goal:** an MCP server that starts on stdio and completes the handshake — no Tiger yet.
- `go mod init`; add `modelcontextprotocol/go-sdk` dependency.
- Directory layout per tech-stack.md (`cmd/tiger-mcp`, `internal/...`).
- `main.go` builds an `mcp.Server` with zero tools and runs `StdioTransport`.
- Logging goes to **stderr**.
- **Verify:** server launches; an MCP client (or the SDK's inspector) sees it connect
  with an empty tool list.

## Phase 1 — Config + Tiger client + one `get_quote` tool  ⬅ first real slice
**Goal:** the thinnest end-to-end path from Claude → server → **live Tiger** → back.
- `internal/config`: load & validate `TIGER_ID` / `TIGER_PRIVATE_KEY` / `TIGER_ACCOUNT`
  from env; fail fast with clear errors.
- Pull the Tiger Go SDK; **read its real source** to confirm client/config/quote APIs.
- `internal/tiger`: construct `ClientConfig` + http client + quote client at startup;
  add a startup smoke check (e.g. market-state) so bad creds fail immediately.
- `internal/tools`: implement one MCP tool `get_quote(symbols)` returning real-time
  quotes for one or more symbols.
- **Verify:** with real creds, ask Claude to quote a symbol and get a live answer.

## Phase 2 — Remaining market-data tools
**Goal:** round out read-only market data. One tool per change if useful.
- `get_market_status(market)`
- `get_kline(symbol, period)`
- `get_depth(symbol)`
- **Verify:** each returns live data; inputs validated with helpful errors.

## Phase 3 — Account & positions tools
**Goal:** portfolio read access (still strictly read-only).
- `get_account_assets()`
- `get_positions()`
- `get_orders()` (history/read only — never place/cancel)
- **Verify:** values reconcile with the Tiger app for the operator's account.

## Phase 4 — Hardening & operator UX
**Goal:** make it pleasant and safe to run daily.
- Consistent tool error handling; validate/sanitize all inputs.
- Sensible timeouts; light rate-limit / retry around Tiger calls.
- Structured stderr logging with levels.
- `README.md`: setup, env vars, and Claude client config snippet (stdio command).
- **Verify:** clean run from a fresh clone following only the README.

## Phase 5 — Tests & polish (as needed)
- Unit tests for config parsing and the tiger wrapper (SDK behind an interface, mocked).
- `go vet` / lint clean; small CI if desired.

---

## Future (out of current scope)
- SSE / Streamable HTTP transport for remote use.
- Sandbox/paper (`TIGER_ENV`) switch.
- Trading tools — **only** behind an explicit, deliberate opt-in with strong guards;
  currently a non-goal (see mission.md).
