package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/muzzarellimj/grace-material-api/database"
	api "github.com/muzzarellimj/grace-material-api/pkg/api/movie"
)

func main() {
	godotenv.Load()

	database.Connect(os.Getenv("DATABASE_CONNECTION_USERNAME"), os.Getenv("DATABASE_CONNECTION_PASSWORD"), os.Getenv("DATABASE_CONNECTION_HOST"), os.Getenv("DATABASE_CONNECTION_PORT"))
	defer database.Disconnect()

	router := gin.Default()

	router.GET("/api/movie", api.HandleGetMovie)

	router.Run("localhost:8080")
}
