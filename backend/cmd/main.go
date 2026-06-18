package main

import (
	"time"

	"citramascoweb-backend/config"
	"citramascoweb-backend/internal/modules/auth"
	"citramascoweb-backend/internal/modules/category"
	"citramascoweb-backend/internal/modules/types"
	"citramascoweb-backend/internal/modules/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.Default()

	db := config.ConnectDB()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}

	app.Use(cors.New(corsConfig))

	port := config.Config("PORT")

	//  init modules
	categoryModule := category.InitModule(db)
	typeModule := types.InitModule(db)
	authModule := auth.InitModule(db)
	userModule := user.InitModule(db)

	api := app.Group("/api")

	categoryModule.CategoryRoutes(api)
	typeModule.TypeRoutes(api)
	authModule.AuthRoutes(api)
	userModule.UserRoutes(api)

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	app.Run(":" + port)
}
