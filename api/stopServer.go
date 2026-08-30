package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// API Request body
type serverStopRequest struct {
	User string `json:"user" binding:"required"`
}

func (server *Server) postStopServer(c *gin.Context) {
	var req serverStartRequest
	var cntID uint64

	// Parameter sent via URL
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepares the command to start the server
	commandStop := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./gameserver stop'"`,
		cntID,
		req.User)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optStartReturn, err := server.sshClient.NewSession(commandStop)

	fmt.Println(optStartReturn)
	fmt.Print(err)
}
