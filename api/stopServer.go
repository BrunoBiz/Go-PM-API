package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *Server) postStopServer(c *gin.Context) {
	var req ServerRequest
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
	optStopReturn, err := server.sshClient.NewSession(commandStop)

	if err != nil {
		slog.Error("SSH New Session: " + err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	// Server stopped
	if strings.Contains(optStopReturn, "[  OK  ] Stopping") ||
		strings.Contains(optStopReturn, "MESSAGE: Server stopped") {
		c.IndentedJSON(http.StatusOK, "Server stopped successfully.")
		return
	}

	// Server is already stopped
	if strings.Contains(optStopReturn, "is already stopped") {
		c.IndentedJSON(http.StatusConflict, "Server is already stopped.")
		return
	}

	// Error
	c.IndentedJSON(http.StatusInternalServerError, "An error occurred, the server could not be stopped.")
}
