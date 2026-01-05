package controllers

import (
	"backend/app/controllers/clientcontrollers"
	"backend/app/middleware"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page     int `json:"page" binding:"min=1"`
	PageSize int `json:"pageSize" binding:"min=1,max=100"`
}

type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

func RegisterControllers(router *gin.RouterGroup) {
	router.POST("/report", middleware.UseClientAuth, middleware.UseGzip, clientcontrollers.ClientController.Report)

	router.POST("/stats", middleware.UseAppAuth, MetricRecordController.FindHomepageStats)

	router.POST("/transactions", middleware.UseAppAuth, TransactionController.FindAllTransactions)
	router.POST("/exception-stack-traces", middleware.UseAppAuth, ExceptionStackTraceController.FindGrouppedExceptionStackTraces)

	// Auth
	router.POST("/login", AuthController.Login)
}
