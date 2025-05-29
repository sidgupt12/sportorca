package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getFixtures(c *gin.Context) {
	apiKey := os.Getenv("API_SPORTS_KEY")
	baseURL := os.Getenv("API_SPORTS_BASE_URL")

	// Get date and status parameters
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	status := c.Query("status")

	// Build query string
	queryString := fmt.Sprintf("date=%s", date)
	if status != "" {
		queryString = fmt.Sprintf("%s&status=%s", queryString, status)
	}

	// Create request URL with query parameters
	url := fmt.Sprintf("%s/fixtures?%s", baseURL, queryString)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}

	// Add headers
	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error making request to API"})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading API response"})
		return
	}

	// Set the content type header
	c.Header("Content-Type", "application/json")
	
	// Write the raw response body directly
	c.Writer.Write(body)
}

func main() {
	loadEnv()

	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Routes
	r.GET("/api/fixtures", getFixtures)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
