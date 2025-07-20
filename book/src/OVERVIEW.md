# Strava MCP Server

A Strava MCP server with local caching for improved performance and rate limit management.

## Prerequisites

1. [Setup Strava developer credentials](./setup-developer-credentials.md)
2. [Install Go](https://go.dev/doc/install)

## Installation

```bash
git clone https://github.com/lyledean1/strava-mcp
cd strava-mcp
make build
```

## Claude Code Configuration

Add to your Claude Code MCP configuration:

```json
{
  "mcpServers": {
    "strava": {
      "command": "{pathToClonedRepo}/bin/strava-mcp",
      "args": [],
      "env": {
        "STRAVA_CLIENT_SECRET": "{clientSecret}",
        "STRAVA_CLIENT_ID": "{clientID}",
        "FOLDER_PATH": "{pathToClonedRepo}"
      }
    }
  }
}
```

Replace `{pathToClonedRepo}`, `{clientSecret}`, and `{clientID}` with your actual values.