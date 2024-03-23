package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	bookApi "github.com/muzzarellimj/grace-material-api/pkg/api/book"
	gameApi "github.com/muzzarellimj/grace-material-api/pkg/api/game"
	movieApi "github.com/muzzarellimj/grace-material-api/pkg/api/movie"
	"github.com/muzzarellimj/grace-material-api/pkg/database/connection"
)

func main() {
	godotenv.Load()

	connection.Connect(os.Getenv("DATABASE_CONNECTION_USERNAME"), os.Getenv("DATABASE_CONNECTION_PASSWORD"), os.Getenv("DATABASE_CONNECTION_HOST"), os.Getenv("DATABASE_CONNECTION_PORT"))

	defer connection.Disconnect()

	router := gin.Default()

	router.GET("/api/book", bookApi.HandleGetBook)
	router.POST("/api/book", bookApi.HandlePostBook)
	router.GET("/api/book/search", bookApi.HandleGetBookSearch)

	router.GET("/api/game", gameApi.HandleGetGame)

	router.GET("/api/movie", movieApi.HandleGetMovie)
	router.POST("/api/movie", movieApi.HandlePostMovie)
	router.GET("/api/movie/search", movieApi.HandleGetMovieSearch)

	router.Run("localhost:8080")
}
