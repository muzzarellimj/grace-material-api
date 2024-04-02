package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	bookApi "github.com/muzzarellimj/grace-material-api/internal/api/book"
	gameApi "github.com/muzzarellimj/grace-material-api/internal/api/game"
	movieApi "github.com/muzzarellimj/grace-material-api/internal/api/movie"
	"github.com/muzzarellimj/grace-material-api/internal/database"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load .env configuration file; this relies on explicitly declared configuration values and may cause issues outside deployed environments: %v\n", err)
	}

	err = database.Connect(os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to persist Grace database pool connection: %v\n", err)

		os.Exit(1)
	}

	defer database.Disconnect()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/api/book", bookApi.HandleGetBook)
	router.POST("/api/book", bookApi.HandlePostBook)
	router.GET("/api/book/search", bookApi.HandleGetBookSearch)

	router.GET("/api/game", gameApi.HandleGetGame)
	router.POST("/api/game", gameApi.HandlePostGame)
	router.GET("/api/game/search", gameApi.HandleGetGameSearch)

	router.GET("/api/movie", movieApi.HandleGetMovie)
	router.POST("/api/movie", movieApi.HandlePostMovie)
	router.GET("/api/movie/search", movieApi.HandleGetMovieSearch)

	err = router.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start listening with router: %v\n", err)
		os.Exit(1)
	}
}
