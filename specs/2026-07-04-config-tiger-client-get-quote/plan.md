# Plan — Config + Tiger client + `get_quote` tool (Phase 1)

Numbered task groups. Each should leave the repo building and (from group 2
onward) with passing tests.

## 1. Verify the real Tiger Go SDK API
- `go get github.com/tigerfintech/openapi-go-sdk@latest`.
- Use `go doc` against the pulled module (same technique as Phase 0's MCP SDK
  check) to confirm actual package names, config constructor, and quote
  client method signature for real-time quotes.
- Do **not** proceed to group 3 until these are confirmed from source —
  earlier signatures from research are unverified (see requirements.md's
  "Open risk").

## 2. Config loading (`internal/config`)
- Define a `Config` struct: `TigerID`, `PrivateKey`, `Account`.
- `Load() (*Config, error)` reads `TIGER_ID`, `TIGER_PRIVATE_KEY`,
  `TIGER_ACCOUNT` from the environment.
- Missing/empty var → error naming the specific missing var (not a generic message).
- Unit tests: all-present succeeds; each var missing individually produces a
  distinct, correct error.

## 3. Tiger client wrapper (`internal/tiger`)
- Define a small interface, e.g.:
  ```go
  type QuoteFetcher interface {
      RealTimeQuotes(ctx context.Context, symbols []string) ([]Quote, error)
  }
  ```
  with a `Quote` struct shaped by what group 1 confirmed the SDK returns.
- Implement it using the real Tiger SDK client (constructed from `config.Config`).
- No unit test for the real implementation itself (would require live
  network/credentials) — but keep it small enough that it's obviously correct
  by inspection, matching the SDK calls confirmed in group 1.

## 4. `get_quote` MCP tool (`internal/tools`)
- Define the tool's input schema (`symbols []string`) and output shape.
- Handler validates input (non-empty `symbols`), calls `QuoteFetcher`, maps
  errors to MCP tool errors, and formats a successful result.
- Unit tests using a fake `QuoteFetcher`: happy path (multiple symbols),
  empty-symbols validation error, upstream-error passthrough.

## 5. Wire-up (`cmd/tiger-mcp/main.go`)
- At startup: `config.Load()` → build the real Tiger client → construct the
  `get_quote` tool → `server.AddTool(...)` → `server.Run(ctx, &mcp.StdioTransport{})`.
- Config load failure → clear fatal error on stderr, non-zero exit, before
  attempting to run the server.

## 6. Tests, vet, and follow-up logging
- `go test ./...`, `go vet ./...`, `gofmt -l .` all clean.
- Record a follow-up note (e.g. in validation.md's results or a TODO) to run
  a **live** MCP Inspector `tools/call get_quote` once real Tiger credentials
  are available — this phase merges without it per the agreed merge bar, but
  the follow-up must not be forgotten.
