package main

import (
	"traceway"
	tracewaygin "traceway/traceway_gin"

	"github.com/gin-gonic/gin"
)

func main() {
	testGin()
}

func testGin() {

	router := gin.Default()

	router.Use(tracewaygin.New(
		"tracewaydemo",
		"default_token_change_me@http://localhost:8082/api/report",
		traceway.WithDebug(true),
	))

	router.GET("/test-exception", func(ctx *gin.Context) {
		panic("Cool")
	})

	router.GET("/test-ok", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.GET("/test-not-found", func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"status": "not-found",
		})
	})

	router.GET("/test-param/:param", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"param": ctx.Param("param"),
		})
	})

	router.GET("/metrics", func(ctx *gin.Context) {
		traceway.PrintCollectionFrameMetrics()
	})

	router.Run()
}
