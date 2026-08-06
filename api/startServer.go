package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type serverStartRequest struct {
	User string `json:"user" binding:"required"`
}

func (server *Server) postStartServer(c *gin.Context) {
	var req serverStartRequest
	var cntID uint64
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commandStart := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./gameserver start'"`,
		cntID,
		req.User)

	optStartReturn, err := server.sshClient.NewSession(commandStart)
	slog.Info("SSH combined output: " + optStartReturn)

	if err != nil {
		slog.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, nil) // TODO - Write a proper return message - Check if already running or started and such -- Use returned message
}
