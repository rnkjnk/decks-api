package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rnkjnk/decks-api/internal/api"
	"github.com/rnkjnk/decks-api/internal/services"
	"github.com/rnkjnk/decks-api/internal/utils"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Load configuration
	config := utils.GetConfigsFromYaml("config.yaml")

	// Initialize repository
	store := services.NewDecksInMemoryStore()

	// Inject dependencies into decks service
	service := services.NewDecksService(config.Decks, store)

	// Inject dependencies into handlers
	handlers := api.NewHandlers(service)

	// Set up routes
	handlers.SetupRoutes(router)

	// Start the server
	router.Run(":" + config.Api.ServerPort)
}
