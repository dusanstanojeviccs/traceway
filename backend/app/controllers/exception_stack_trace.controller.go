package controllers

import (
	"backend/app/models"
	"backend/app/repositories"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type exceptionStackTraceController struct{}

type ExceptionSearchRequest struct {
	FromDate   time.Time        `json:"fromDate"`
	ToDate     time.Time        `json:"toDate"`
	OrderBy    string           `json:"orderBy"`
	Pagination PaginationParams `json:"pagination"`
	Search     string           `json:"search"`
}

func (e exceptionStackTraceController) FindGrouppedExceptionStackTraces(c *gin.Context) {
	var request ExceptionSearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exceptions, total, err := repositories.ExceptionStackTraceRepository.FindGrouped(c, request.FromDate, request.ToDate, request.Pagination.Page, request.Pagination.PageSize, request.OrderBy, request.Search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse[models.ExceptionGroup]{
		Data: exceptions,
		Pagination: Pagination{
			Page:       request.Pagination.Page,
			PageSize:   request.Pagination.PageSize,
			Total:      total,
			TotalPages: (total + int64(request.Pagination.PageSize) - 1) / int64(request.Pagination.PageSize),
		},
	})
}

var ExceptionStackTraceController = exceptionStackTraceController{}
