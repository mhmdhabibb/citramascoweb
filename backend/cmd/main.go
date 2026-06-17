package main

import (
	"citramascoweb-backend/config"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.Default()

	config.ConnectDB()

	db := config.Config("DB_URL")
	port := config.Config("PORT")

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}

	app.Use(cors.New(corsConfig))
	fmt.Println("Database: ", db)

	// api := app.Group("/api")

	app.Run(":", port)
}
