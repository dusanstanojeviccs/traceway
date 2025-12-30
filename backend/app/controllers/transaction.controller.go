package controllers

import (
	"backend/app/models"
	"backend/app/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type transactionController struct{}

type TransactionSearchRequest struct {
	FromDate   time.Time        `json:"fromDate"`
	ToDate     time.Time        `json:"toDate"`
	OrderBy    string           `json:"orderBy"`
	Pagination PaginationParams `json:"pagination"`
}

func (e transactionController) FindAllTransactions(c *gin.Context) {
	var request TransactionSearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, total, err := repositories.TransactionRepository.FindAll(c, request.FromDate, request.ToDate, request.Pagination.Page, request.Pagination.PageSize, request.OrderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse[models.Transaction]{
		Data: transactions,
		Pagination: Pagination{
			Page:       request.Pagination.Page,
			PageSize:   request.Pagination.PageSize,
			Total:      total,
			TotalPages: (total + int64(request.Pagination.PageSize) - 1) / int64(request.Pagination.PageSize),
		},
	})
}

var TransactionController = transactionController{}
