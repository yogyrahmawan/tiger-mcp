# Validation — Repo scaffold & empty server (Phase 0)

How to know this phase is done and safe to merge.

## Build
- [ ] `go build ./...` succeeds from a clean checkout of this branch.
- [ ] `go vet ./...` reports no issues.

## Runtime shape
- [ ] Running the built binary directly does not crash, panic, or print anything
      to **stdout** (stdout must stay clean for the MCP protocol stream).
- [ ] Any startup logs appear on **stderr** only.
- [ ] The process blocks waiting for stdio input rather than exiting immediately
      (confirms it's actually running the server loop, not a no-op).

## Protocol-level check (MCP Inspector)
- [ ] Launch the binary under `@modelcontextprotocol/inspector` (or equivalent
      official inspector tooling) configured for stdio transport.
- [ ] Inspector reports a successful `initialize` handshake.
- [ ] Inspector's tool list for this server is **empty** (expected — no tools
      exist yet in Phase 0).
- [ ] No unexpected errors/warnings in the inspector session.

## Structure check
- [ ] Module path in `go.mod` is `github.com/yogyrahmawan/tiger-mcp`.
- [ ] Directory layout matches [tech-stack.md](../tech-stack.md#project-layout-target):
      `cmd/tiger-mcp/`, `internal/config/`, `internal/tiger/`, `internal/tools/`.
- [ ] No Tiger SDK dependency, no credential handling, no tool implementations
      present yet (would indicate scope creep into Phase 1).

## Definition of done
This phase is mergeable when every checkbox above is checked, the branch
`phase-0-repo-scaffold` builds cleanly, and the inspector run has been
performed at least once by hand (paste/note the command + outcome in the PR
description).
