package middleware

import (
	"backend/app/db"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTx(c *gin.Context) *sql.Tx {
	// this will panic if we call it in a non transactional route - we're ok with that
	txFromC, exists := c.Get("tx")
	if !exists {
		panic("No TX bound to the context, please make sure you're using the UseTransaction middleware")
	}

	tx, ok := txFromC.(*sql.Tx)

	if !ok {
		panic("No TX bound to the context, please make sure you're using the UseTransaction middleware")
	}

	return tx
}

// only commits the tx for status 200 - ignores 201
func UseTransaction(c *gin.Context) {
	tx, err := db.DB.Begin()
	if err != nil {
		_ = c.Error(fmt.Errorf("failed to begin transaction: %w", err))
		c.Abort()
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if status := c.Writer.Status(); status == http.StatusOK || status == http.StatusCreated {
			if err := tx.Commit(); err != nil {
				_ = c.Error(fmt.Errorf("failed to commit: %w", err))
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
	}()

	c.Set("tx", tx)
	c.Next()
}
