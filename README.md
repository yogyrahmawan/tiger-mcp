# tiger-mcp

An MCP (Model Context Protocol) server exposing read-only Tiger Open API
market data and account access to MCP clients. See [specs/mission.md](specs/mission.md),
[specs/tech-stack.md](specs/tech-stack.md), and [specs/roadmap.md](specs/roadmap.md)
for the project's constitution and implementation plan.

## Build & run

```sh
go build -o tiger-mcp ./cmd/tiger-mcp
./tiger-mcp
```

The server communicates over stdio using the MCP protocol. As of Phase 0 it
registers no tools yet — see [specs/roadmap.md](specs/roadmap.md) for what's next.

## Input from stakeholder
- I want to build tigertrade MCP with golang 