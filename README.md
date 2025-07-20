# Strava MCP Server

A Model Context Protocol (MCP) server for accessing Strava activity data through Claude, with local caching for improved performance.

## Quick Start

1. [Setup Strava developer credentials](./setup-developer-credentials.md)
2. Install Go and clone the repository:
   ```bash
   git clone https://github.com/lyledean1/strava-mcp
   cd strava-mcp
   make build
   ```
3. Configure Claude Code with your credentials:
   ```json
   {
     "mcpServers": {
       "strava": {
         "command": "{pathToClonedRepo}/bin/strava-mcp",
         "env": {
           "STRAVA_CLIENT_SECRET": "{clientSecret}",
           "STRAVA_CLIENT_ID": "{clientID}",
           "FOLDER_PATH": "{pathToClonedRepo}"
         }
       }
     }
   }
   ```

## What You Can Do

- **Analyze Activities**: Get detailed stats on runs, rides, swims, and other activities
- **Filter by Date/Type**: Find specific activities with flexible filtering
- **Stream Data**: Access GPS, heart rate, power, and sensor data for detailed analysis
- **Performance Insights**: Compare workouts, track progress, and identify patterns

## Features

- Local data caching for fast responses
- Comprehensive activity metrics (distance, speed, heart rate, power, elevation)
- Time-series stream data for detailed analysis
- Automatic rate limit management

## Documentation

This project is part of a comprehensive guide to building MCP servers. For detailed documentation and examples, see the accompanying book.

## Tools Available

- `get_activities` - Retrieve and filter your Strava activities
- `get_activity_stream` - Get detailed sensor data for specific activities

Ask Claude to help analyze your fitness data, create visualizations, or track your training progress!

## License
This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.
