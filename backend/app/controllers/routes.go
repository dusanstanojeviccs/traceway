package controllers

import (
	"backend/app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(router *gin.RouterGroup) {
	router.GET("/exceptions", middleware.UseTransaction, ExceptionController.FindAll)
}
