package main

import (
	"log"
	"stravamcp/api"
	"stravamcp/config"
	"stravamcp/pkg/client"
	"stravamcp/repo"
	"stravamcp/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Unable to get config %s", err)
	}
	stravaClient := client.NewStravaClient("https://www.strava.com")
	tokenRepo := repo.NewTokenRepo(stravaClient, cfg.StravaClientID, cfg.StravaClientSecret, cfg.FolderPath, cfg.RefreshTokenFileName)
	activityService := service.NewActivityService(stravaClient, tokenRepo, repo.NewStorage(cfg.FolderPath))
	server := api.SetupRouter(activityService)
	err = server.Run("localhost:8081")
	if err != nil {
		log.Fatalf("Unable to start server %s", err)
	}
}
