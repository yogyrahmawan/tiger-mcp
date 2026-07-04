# Tech Stack

## Language & toolchain
- **Go** (module-based). Target the Go version required by the Tiger SDK (>= 1.20)
  or newer; pin the exact version in `go.mod`.
- Standard `go build` / `go test`. No codegen or heavy frameworks.

## Core dependencies
| Concern | Choice | Import path |
| --- | --- | --- |
| MCP server | **Official Go SDK** (maintained with Google) | `github.com/modelcontextprotocol/go-sdk/mcp` |
| Tiger Open API | **Official Go SDK** | `github.com/tigerfintech/openapi-go-sdk` (`config`, `client`, `quote`, `trade`) |

> Note: exact Tiger SDK package/method names must be **verified against the pulled
> module source** before use — do not trust secondhand signatures.

## Transport
- **stdio only** for now. The server runs as a local subprocess launched by the
  Claude client. (`mcp.StdioTransport`.) SSE / Streamable HTTP are out of scope
  until a remote use case appears.

## Configuration & secrets
- Credentials come from **environment variables**, read at startup:
  - `TIGER_ID` — Tiger developer / tiger id
  - `TIGER_PRIVATE_KEY` — RSA private key (PEM; supplied as env value or `@/path` to a file, TBD in impl)
  - `TIGER_ACCOUNT` — trading account id
- Nothing is written to disk by the server. No secrets in the repo, no properties
  file committed. Missing/blank values cause an immediate, descriptive startup error.

## Environment
- Targets the operator's **live Tiger account**. There is no sandbox toggle in
  scope; because the server is read-only, live is acceptable. (A `TIGER_ENV`
  switch can be added later without changing tool shapes.)

## Project layout (target)
```
tiger-mcp/
  cmd/tiger-mcp/main.go    # entrypoint: load config, build clients, run stdio server
  internal/config/         # env loading + validation
  internal/tiger/          # thin wrapper around Tiger SDK clients (quote, account)
  internal/tools/          # MCP tool definitions + handlers
  specs/                   # this constitution
  README.md
```

## Auth model (Tiger)
- Tiger Open API uses **RSA signing**. The private key + tiger id + account are
  loaded once at startup into a Tiger `ClientConfig`; a single HTTP client and the
  quote/trade clients are constructed from it and reused.

## Testing & quality
- Unit-test config parsing/validation and the tiger wrapper (with the SDK client
  mocked behind a small interface).
- Manual/live smoke test for the first real Tiger call.
- `go vet` clean; simple structured logging to **stderr** (stdout is reserved for
  the MCP protocol on stdio).

## Deliberately excluded
- No database, no persistence layer.
- No trading endpoints wired up.
- No web framework, no HTTP server.
