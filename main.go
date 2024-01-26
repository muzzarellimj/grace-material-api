package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/example"
)

func main() {
	router := gin.Default()

	// EXAMPLE
	router.GET("/api/ex/games", example.FetchGames)
	router.GET("/api/ex/games/:id", example.FetchGame)
	router.GET("/api/ex/movies", example.FetchMovies)
	router.GET("/api/ex/movies/:id", example.FetchMovie)

	router.Run("localhost:8080")
}
