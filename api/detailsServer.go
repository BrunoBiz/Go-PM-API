package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *Server) postDetailsServer(c *gin.Context) {
	var req ServerRequest
	var cntID uint64

	// Parameter sent via URL
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepares the command to check server details
	commandDetails := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./gameserver details'"`,
		cntID,
		req.User)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optDetailsReturn, err := server.sshClient.NewSession(commandDetails)

	if err != nil {
		slog.Error("SSH New Session: " + err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	// Server Details
	if strings.Contains(strings.ToUpper(optDetailsReturn), "STATUS:") {

	}

	// Error
}
