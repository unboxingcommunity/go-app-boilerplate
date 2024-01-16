package middleware

import (
	"errors"
	"fmt"
	"go-boilerplate-api/apm"
	log "go-boilerplate-api/pkg/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	pkgErrors "github.com/pkg/errors"
	"stash.bms.bz/bms/monitoringsystem"
)

// HandlePanic ... rest panic handler
func HandlePanic(c *gin.Context) {
	defer func(c *gin.Context) {
		r := recover()
		var stackTrace string
		if r != nil {
			err, ok := r.(error)
			if ok {
				// Logs the error
				stackTrace = fmt.Sprintf("%+v", pkgErrors.New(err.Error()))
				log.Error("GO-BOILERPLATE.PANIC", "Unexpected panic occured", log.Priority1, nil, map[string]interface{}{"error": err.Error(), "stackTrace": stackTrace})

				// Notice error in apm
				apm.APM.NoticeError(apm.FromContext(c), err)

				// Forms error message
				c.JSON(500, gin.H{
					"message":    "Panic: Unexpected error occured.",
					"error":      err.Error(),
					"stackTrace": stackTrace,
				})
			} else {
				// Logs the error
				log.Error("GO-BOILERPLATE.PANIC", "Panic recovery failed to parse error", log.Priority1, nil, map[string]interface{}{"error": r})

				// Notice error in apm
				apm.APM.NoticeError(apm.FromContext(c), errors.New("GO-BOILERPLATE.UNRECOVERED.PANIC"))

				// Forms error message
				c.JSON(500, gin.H{
					"message": "Panic: Unexpected error occured, failed to parse error",
				})
			}
		}
	}(c)
	c.Next()
}

// ApmMiddleware creates a middleware to start and end apm transaction.
//
//	router := gin.Default()
//	// Add the the middleware before other middlewares or routes:
//	router.Use(Middleware(apm))
//
// Input
//		apmAgent: instance of the apm
// Output
//		ginFun: gin http handler function
func ApmMiddleware(apmAgent *monitoringsystem.Agent) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apmAgent != nil {
			// Starts an APM transaction
			name := c.HandlerName()

			// Removes any sensitive query params here ... (skip / remove block if no such params exist)
			req := c.Request
			rawQuery := req.URL.RawQuery
			queryParams := req.URL.Query()
			// queryParams.Del("sensitive_query_params")
			req.URL.RawQuery = queryParams.Encode()

			txn, err := apmAgent.StartWebTransaction(name, c.Writer, req)
			c.Request.URL.RawQuery = rawQuery

			// defer End transaction
			defer func(c *gin.Context) {
				var err error
				if c.Writer.Status() != http.StatusOK {
					err = errors.New("GO-BOILERPLATE.ERROR")
				}
				apmAgent.EndTransaction(apm.FromContext(c), err)
			}(c)

			// Handles error incase transaction fails
			if err != nil {
				// Place holder code until new errors library is implemented properly
				log.Error("GO-BOILERPLATE.REST.APM_TRANS_INIT_FAIL", "Transaction failed", log.Priority1, nil, map[string]interface{}{"error": err.Error()})
				c.Next()
			} else if txn == nil {
				log.Error("GO-BOILERPLATE.REST.APM_TRANS_INIT_FAIL", "Transaction is nil", log.Priority1, nil, map[string]interface{}{"transaction": txn})
				c.Next()
			}

			// Stores transaction details in context
			c.Set(apm.TransactionKey, txn)
		}
		c.Next()
	}
}
