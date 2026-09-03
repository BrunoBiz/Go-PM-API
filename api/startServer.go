package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *Server) postStartServer(c *gin.Context) {
	var req ServerRequest
	var cntID uint64

	// Parameter sent via URL
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepares the command to start the server
	commandStart := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./gameserver start'"`,
		cntID,
		req.User)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optStartReturn, err := server.sshClient.NewSession(commandStart)

	if err != nil {
		slog.Error("SSH New Session: " + err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	// Started
	if strings.Contains(optStartReturn, "[  OK  ] Starting") ||
		strings.Contains(optStartReturn, "MESSAGE: Server started") {
		c.IndentedJSON(http.StatusOK, "Server started successfully.")
		return
	}

	// Already running
	if strings.Contains(optStartReturn, "is already running") {
		c.IndentedJSON(http.StatusConflict, "Server is already running.")
		return
	}

	// Error
	c.IndentedJSON(http.StatusInternalServerError, "An error has occurred, the server could not be started.")
}
