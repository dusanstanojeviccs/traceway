package tracewaygin

import (
	"time"
	"traceway"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func wrapAndExecute(c *gin.Context) (s *string) {
	defer func() {
		if r := recover(); r != nil {
			m := traceway.FormatRWithStack(r, traceway.CaptureStack(2))
			s = &m
		}
	}()
	c.Next()
	return nil
}

func New(app, connectionString string, options ...func(*traceway.TracewayOptions)) gin.HandlerFunc {
	traceway.Init(app, connectionString, options...)

	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()
		transactionId := uuid.NewString()

		stackTraceFormatted := wrapAndExecute(c)

		duration := time.Since(start)

		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if query != "" {
			path = path + "?" + query
		}

		transactionEndpoint := method + " " + path

		defer recover()

		traceway.CaptureTransaction(transactionId, transactionEndpoint, duration, start, statusCode, bodySize, clientIP)

		if stackTraceFormatted != nil {
			traceway.CaptureTransactionException(transactionId, *stackTraceFormatted)
		}
	}
}
