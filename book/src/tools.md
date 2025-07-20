# Strava MCP Server Tools

This MCP server provides tools to interact with your Strava activity data through Claude. All data is cached locally for improved performance and reduced API rate limits.

## Available Tools [WIP]

### `get_activities`
Retrieve your Strava activities with optional filtering and date range support.

**Parameters:**
- `filter` (optional): Activity type filter (e.g., 'runs', 'rides', 'swims')
- `before` (optional): Return activities before this date (ISO 8601 format)
- `after` (optional): Return activities after this date (ISO 8601 format)

**Returns:**
- Activity summary with total count and applied filters
- Detailed activity information including:
    - Activity name and ID
    - Activity type (run, ride, swim, etc.)
    - Distance (in kilometers)
    - Duration (hours, minutes, seconds)
    - Average and max heart rate
    - Average and max speed
    - Power data (average and weighted average watts)
    - Elevation gain
    - Start date and location coordinates

**Example Usage:**
```
Get all my runs from last month
Get cycling activities after 2024-01-01
Show activities before 2024-06-01T00:00:00Z
```

### `get_activity_stream`
Get detailed stream data for a specific activity including GPS coordinates, heart rate, power, and other sensor data.

**Parameters:**
- `activity_id` (required): The ID of the activity to retrieve stream data for

**Returns:**
- Activity name and ID
- Number of data points collected
- Available data types (time, power, heart rate, etc.)
- Complete stream data with time-series information

**Example Usage:**
```
Get stream data for activity 12345678
Show detailed GPS and heart rate data for my latest run
```

## Data Format

Activities include comprehensive metrics when available:
- **Basic Info**: Name, type, date
- **Performance**: Distance, speed, duration
- **Health**: Heart rate data
- **Power**: Watts for cycling activities
- **Geography**: Elevation gain, GPS coordinates

Stream data provides time-series information for detailed analysis and visualization of your workout data.

## Notes

- All dates should be in ISO 8601 format (e.g., `2024-01-15T00:00:00Z`)
- Data is automatically cached locally to minimize API calls
- The server respects Strava's rate limits through intelligent caching