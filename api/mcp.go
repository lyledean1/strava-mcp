package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"stravamcp/service"
	"strings"
	"time"
)

type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MCPServer struct {
	activityService service.ActivityService
	upgrader        websocket.Upgrader
}

func NewMCPServer(activityService service.ActivityService) *MCPServer {
	return &MCPServer{
		activityService: activityService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

func (s *MCPServer) HandleMCPRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return s.getMcpInitResponse(req)

	case "tools/list":
		return s.getToolList(req)

	case "tools/call":
		return s.handleToolCall(req)

	case "notifications/initialized":
		return MCPResponse{}

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func (s *MCPServer) getMcpInitResponse(req MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{
					"listChanged": false,
				},
			},
			"serverInfo": map[string]interface{}{
				"name":    "strava-mcp",
				"version": "1.0.0",
			},
		},
	}
}

func (s *MCPServer) getToolList(req MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": []map[string]interface{}{
				{
					"name":        "get_activities",
					"description": "Get all activities with optional filtering (runs, rides, swims, etc.)",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"filter": map[string]interface{}{
								"type":        "string",
								"description": "Activity type filter (e.g., 'runs', 'rides', 'swims')",
							},
							"before": map[string]interface{}{
								"type":        "string",
								"description": "Return activities before this date (ISO 8601 format)",
							},
							"after": map[string]interface{}{
								"type":        "string",
								"description": "Return activities after this date (ISO 8601 format)",
							},
						},
					},
				},
				{
					"name":        "get_activity_stream",
					"description": "Get detailed stream data for a specific activity (GPS, heart rate, power, etc.)",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"activity_id": map[string]interface{}{
								"type":        "string",
								"description": "The ID of the activity",
							},
						},
						"required": []string{"activity_id"},
					},
				},
				// comment this out for the time being, don't want it to be manually refreshing
				//{
				//	"name":        "refresh_activities",
				//	"description": "Refresh activities from Strava API to get latest data",
				//	"inputSchema": map[string]interface{}{
				//		"type": "object",
				//	},
				//},
			},
		},
	}
}

func (s *MCPServer) handleToolCall(req MCPRequest) MCPResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Missing tool name",
			},
		}
	}

	arguments, _ := params["arguments"].(map[string]interface{})

	c := &gin.Context{}

	switch toolName {
	case "get_activities":
		return s.getActivities(req, arguments, c)

	case "get_activity_stream":
		return s.getActivityStream(req, arguments, c)

	case "refresh_activities":
		return s.refreshActivities(req)

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Tool not found",
			},
		}
	}
}

func (s *MCPServer) getActivities(req MCPRequest, arguments map[string]interface{}, c *gin.Context) MCPResponse {
	filter := ""
	if f, ok := arguments["filter"].(string); ok {
		filter = f
	}

	var before, after *time.Time

	if b, ok := arguments["before"].(string); ok && b != "" {
		if parsed, err := time.Parse(time.RFC3339, b); err == nil {
			before = &parsed
		} else {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &MCPError{
					Code:    -32602,
					Message: "Invalid 'before' date format. Use ISO 8601 format (e.g., 2024-01-15T00:00:00Z)",
				},
			}
		}
	}

	// Parse after date
	if a, ok := arguments["after"].(string); ok && a != "" {
		if parsed, err := time.Parse(time.RFC3339, a); err == nil {
			after = &parsed
		} else {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &MCPError{
					Code:    -32602,
					Message: "Invalid 'after' date format. Use ISO 8601 format (e.g., 2024-01-15T00:00:00Z)",
				},
			}
		}
	}

	// Validate date range
	if before != nil && after != nil && before.Before(*after) {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "'before' date must be after 'after' date",
			},
		}
	}

	activities, err := s.activityService.GetAllActivities(c, filter, before, after)
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: fmt.Sprintf("Failed to get activities: %v", err),
			},
		}
	}

	// Create descriptive summary
	summary := fmt.Sprintf("Retrieved %d activities", len(activities))

	var filters []string
	if filter != "" {
		filters = append(filters, fmt.Sprintf("type: %s", filter))
	}
	if after != nil {
		filters = append(filters, fmt.Sprintf("after: %s", after.Format("2006-01-02")))
	}
	if before != nil {
		filters = append(filters, fmt.Sprintf("before: %s", before.Format("2006-01-02")))
	}

	if len(filters) > 0 {
		summary += fmt.Sprintf(" (filtered by %s)", strings.Join(filters, ", "))
	}

	// Format activities for display
	var contentItems []map[string]interface{}

	// Add summary as first item
	contentItems = append(contentItems, map[string]interface{}{
		"type": "text",
		"text": summary,
	})

	// Add each activity as a formatted text item
	for _, activity := range activities {
		var activityText strings.Builder

		// Format the activity details in a readable way
		activityText.WriteString(fmt.Sprintf("🏃 %s (ID: %d)\n", activity.Name, activity.ID))

		if activity.Type != "" {
			activityText.WriteString(fmt.Sprintf("   Type: %s\n", activity.Type))
		}

		if activity.Distance > 0 {
			activityText.WriteString(fmt.Sprintf("   Distance: %.2f km\n", activity.Distance/1000))
		}

		if activity.MovingTime > 0 {
			hours := activity.MovingTime / 3600
			minutes := (activity.MovingTime % 3600) / 60
			seconds := activity.MovingTime % 60
			if hours > 0 {
				activityText.WriteString(fmt.Sprintf("   Duration: %dh %dm %ds\n", hours, minutes, seconds))
			} else {
				activityText.WriteString(fmt.Sprintf("   Duration: %dm %ds\n", minutes, seconds))
			}
		}

		if activity.AverageHeartrate != nil {
			activityText.WriteString(fmt.Sprintf("   Avg heart rate: %.2f bpm\n", *activity.AverageHeartrate))
		}

		if activity.AverageSpeed > 0 {
			activityText.WriteString(fmt.Sprintf("   Avg Speed: %.2f km/h\n", activity.AverageSpeed*3.6))
		}

		if activity.MaxSpeed > 0 {
			activityText.WriteString(fmt.Sprintf("   Max Speed: %.2f km/h\n", activity.MaxSpeed*3.6))
		}

		if activity.AverageWatts != nil {
			activityText.WriteString(fmt.Sprintf("   Average Watts (Power): %.2f \n", *activity.AverageWatts))
		}

		if activity.WeightedAverageWatts != nil {
			activityText.WriteString(fmt.Sprintf("   Weighted Average Watts (Power): %d \n", *activity.WeightedAverageWatts))
		}

		if activity.TotalElevationGain > 0 {
			activityText.WriteString(fmt.Sprintf("   Elevation Gain: %.0f m\n", activity.TotalElevationGain))
		}

		if activity.StartDate != "" {
			activityText.WriteString(fmt.Sprintf("   Date: %s\n", activity.StartDate))
		}

		if len(activity.StartLatLng) >= 2 {
			activityText.WriteString(fmt.Sprintf("   Start Location: %.6f, %.6f\n", activity.StartLatLng[0], activity.StartLatLng[1]))
		}

		contentItems = append(contentItems, map[string]interface{}{
			"type": "text",
			"text": activityText.String(),
		})
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": contentItems,
			"data":    activities, // Keep the raw data for programmatic access if needed
		},
	}
}

func (s *MCPServer) getActivityStream(req MCPRequest, arguments map[string]interface{}, c *gin.Context) MCPResponse {
	activityID, ok := arguments["activity_id"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Missing or invalid activity_id parameter",
			},
		}
	}

	activityStream, err := s.activityService.GetActivityStream(c, activityID)
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: fmt.Sprintf("Failed to retrieve activity stream data: %v", err),
			},
		}
	}

	streamCount := len(activityStream.Streams)
	summaryText := fmt.Sprintf("Retrieved stream data for activity %s", activityID)

	if activityStream.ActivityName != "" {
		summaryText = fmt.Sprintf("Retrieved stream data for '%s' (ID: %s)",
			activityStream.ActivityName, activityID)
	}

	summaryText += fmt.Sprintf("\n- %d data points collected", streamCount)

	// Add details about available data types
	if streamCount > 0 {
		var dataTypes []string
		firstPoint := activityStream.Streams[0]
		if firstPoint.Time != nil {
			dataTypes = append(dataTypes, "time")
		}
		if firstPoint.Watts != nil {
			dataTypes = append(dataTypes, "power")
		}
		if firstPoint.Heartrate != nil {
			dataTypes = append(dataTypes, "heart rate")
		}

		if len(dataTypes) > 0 {
			summaryText += fmt.Sprintf("\n- Available data: %s", strings.Join(dataTypes, ", "))
		}
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": summaryText,
				},
			},
			"stream_data": activityStream,
		},
	}
}

func (s *MCPServer) refreshActivities(req MCPRequest) MCPResponse {
	err := s.activityService.ProcessActivities()
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: fmt.Sprintf("Failed to refresh activities: %v", err),
			},
		}
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "Activities refreshed successfully from Strava API",
				},
			},
		},
	}
}
