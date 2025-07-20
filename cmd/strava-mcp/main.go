package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"stravamcp/api"
	"stravamcp/config"
	"stravamcp/pkg/client"
	"stravamcp/repo"
	"stravamcp/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	stravaClient := client.NewStravaClient("https://www.strava.com")
	tokenRepo := repo.NewTokenRepo(stravaClient, cfg.StravaClientID, cfg.StravaClientSecret, cfg.FolderPath, cfg.RefreshTokenFileName)
	activityService := service.NewActivityService(stravaClient, tokenRepo, repo.NewStorage(cfg.FolderPath))
	mcpServer := api.NewMCPServer(activityService)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		fmt.Fprintf(os.Stderr, "Received: %s\n", line)

		var request api.MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			fmt.Fprintf(os.Stderr, "JSON parse error: %v\n", err)
			// Send a proper error response
			errorResponse := api.MCPResponse{
				JSONRPC: "2.0",
				ID:      nil,
				Error: &api.MCPError{
					Code:    -32700,
					Message: "Parse error",
				},
			}
			responseJSON, _ := json.Marshal(errorResponse)
			fmt.Println(string(responseJSON))
			continue
		}

		fmt.Fprintf(os.Stderr, "Parsed request - Method: %s, ID: %v\n", request.Method, request.ID)

		response := mcpServer.HandleMCPRequest(request)

		if request.Method == "notifications/initialized" || (response.JSONRPC == "" && response.ID == nil && response.Result == nil && response.Error == nil) {
			fmt.Fprintf(os.Stderr, "Skipping response for notification: %s\n", request.Method)
			continue
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Response marshal error: %v\n", err)
			continue
		}

		fmt.Fprintf(os.Stderr, "Sending: %s\n", string(responseJSON))
		fmt.Println(string(responseJSON))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Scanner error: %v\n", err)
	}
}
