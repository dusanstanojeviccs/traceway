package controllers

import (
	"backend/app/middleware"
	"backend/app/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type exceptionController struct{}

func (e exceptionController) FindAll(c *gin.Context) {
	tx := middleware.GetTx(c)
	exceptions, err := repositories.ExceptionRepository.FindAll(tx)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"exceptions": exceptions,
	})
}

var ExceptionController = exceptionController{}
