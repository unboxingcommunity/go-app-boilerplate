package ping

import (
	"go-boilerplate-api/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

// pingResponse is a struct storing the response.
// Message : the message field of the response.
type pingResponse struct {
	Message string `json:"message"`
	Verison string `json:"version"`
}

// pongResponse is a constant used for sending response message.
const pongResponse string = "PONG"

// Ping ... used as a pointer receiver
type Ping struct {
}

// NewPingService ...
func NewPingService() *Ping {
	return &Ping{}
}

// Get returns a response for the /ping request.
// ctx *gin.Context : allows us to pass variables between middleware
func (ping *Ping) Get(ctx *gin.Context) {
	response := pingResponse{Message: pongResponse, Verison: shared.VERSION}
	ctx.JSON(http.StatusOK, response)
}
