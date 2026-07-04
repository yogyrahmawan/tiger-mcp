# Requirements — Repo scaffold & empty server (Phase 0)

## Scope
Stand up the Go module and project skeleton, and get a **bare MCP server**
(zero tools) running over stdio and completing the MCP protocol handshake.
No Tiger integration, no config loading, no tools yet — that is Phase 1.

This is the roadmap's [Phase 0](../roadmap.md#phase-0--repo-scaffold--empty-server).

## In scope
- `go.mod` for module `github.com/yogyrahmawan/tiger-mcp`.
- Dependency on the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk`).
- Directory layout per [tech-stack.md](../tech-stack.md#project-layout-target):
  `cmd/tiger-mcp/`, `internal/config/`, `internal/tiger/`, `internal/tools/`
  (the latter three may be empty/placeholder packages for now).
- `cmd/tiger-mcp/main.go`: constructs an `mcp.Server` with **no tools registered**
  and runs it over `mcp.StdioTransport`.
- Logging to **stderr** only (stdout is reserved for the MCP protocol stream).
- The binary builds and runs without crashing or hanging.

## Out of scope (deferred to later phases)
- Any Tiger SDK dependency, client, or credentials (Phase 1).
- Any MCP tool implementation, including `get_quote` (Phase 1).
- Config loading/validation from env vars (Phase 1).
- Tests beyond a manual handshake check (Phase 5).

## Decisions
- **Module path:** `github.com/yogyrahmawan/tiger-mcp` — matches the existing
  `origin` git remote (`git@github.com:yogyrahmawan/tiger-mcp.git`), so `go get`
  and import paths resolve correctly if/when this is pushed.
- **Go version:** pin `go.mod` to the toolchain version actually installed in
  this environment (verify with `go version` during implementation), as long as
  it satisfies the SDK's minimum (>= 1.20 per tech-stack.md).
- **Task granularity:** coarse — Phase 0 is small enough that a handful of
  numbered task groups is clearer than many micro-steps.
- **Validation method:** use the official `@modelcontextprotocol/inspector`
  tool against the built binary to get real protocol-level proof of a
  successful handshake with an empty tool list, rather than just `go build`/`go vet`
  or wiring into a live Claude client config (saved for later phases where
  there's actually a tool worth exercising end-to-end).

## Context
- Constitution docs: [mission.md](../mission.md), [tech-stack.md](../tech-stack.md),
  [roadmap.md](../roadmap.md).
- Working branch: `phase-0-repo-scaffold`.
- This phase exists purely to de-risk the MCP plumbing before Tiger's live
  account and RSA auth enter the picture in Phase 1.
