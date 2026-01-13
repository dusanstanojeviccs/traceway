package main

import (
	"backend/app/cache"
	"backend/app/chdb"
	"backend/app/controllers"
	"backend/app/middleware"
	"backend/app/migrations"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// we don't actually care for the .env file existing
		// because in production we can just deploy with container variables
		log.Println("Error loading .env file")
	}

	err = chdb.Init()
	if err != nil {
		panic(err)
	}

	err = migrations.Run()
	if err != nil {
		panic(err)
	}

	// Initialize project cache
	ctx := context.Background()
	if err := cache.ProjectCache.Init(ctx); err != nil {
		panic(err)
	}

	middleware.InitUseClientAuth()

	router := gin.Default()

	router.Use(gin.Recovery())

	apiRouterGroup := router.Group("/api")
	controllers.RegisterControllers(apiRouterGroup)

	router.GET("/version", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"version": "0.0.1"})
	})

	if err := router.Run(":8082"); err != nil {
		panic(err)
	}
}
