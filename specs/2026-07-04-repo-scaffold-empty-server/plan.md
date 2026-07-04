# Plan — Repo scaffold & empty server (Phase 0)

Coarse, numbered task groups. Each group should leave the repo in a working
(building) state.

## 1. Module & dependency init
- Confirm installed Go toolchain version (`go version`).
- `go mod init github.com/yogyrahmawan/tiger-mcp`.
- `go get github.com/modelcontextprotocol/go-sdk` (mcp package), pinning to the
  latest tagged release.
- Commit `go.mod` / `go.sum`.

## 2. Directory scaffold
- Create `cmd/tiger-mcp/`, `internal/config/`, `internal/tiger/`, `internal/tools/`
  per [tech-stack.md](../tech-stack.md#project-layout-target).
- Add minimal placeholder (e.g. a `doc.go` or package comment) in each empty
  `internal/*` package so they compile as valid, intentional Go packages rather
  than dangling empty directories.

## 3. Empty stdio server
- Implement `cmd/tiger-mcp/main.go`:
  - Build `mcp.NewServer(&mcp.Implementation{Name: "tiger-mcp", Version: "0.0.1"}, nil)`
    with **no tools registered**.
  - Run it via `server.Run(context.Background(), &mcp.StdioTransport{})`.
  - Any startup logging goes to `os.Stderr` (never stdout).
  - On fatal error, log to stderr and exit non-zero.
- `go build ./...` succeeds; binary starts and blocks waiting on stdio (doesn't
  crash immediately).

## 4. Verify handshake
- Run the built `tiger-mcp` binary under the `@modelcontextprotocol/inspector`
  CLI/UI.
- Confirm: inspector connects, completes the MCP `initialize` handshake, and
  lists **zero tools** (as expected for this phase).
- Capture the exact command used in `validation.md` results / README note for
  reuse in later phases.

## 5. Wrap-up
- `go vet ./...` clean.
- Update root `README.md` with a one-line "build & run" note if not already
  covered elsewhere (kept minimal — full operator docs are Phase 4).
- Open PR / mark branch ready per your usual review flow.
