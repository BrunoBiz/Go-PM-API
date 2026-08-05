package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type serverOptRequest struct {
	LxcId       int    `json:"lxcId" binding:"required"`
	User        string `json:"user" binding:"required"`
	Placeholder string `json:"Placeholder" binding:"required"` // TODO
}

func (server *Server) postStartServer(c *gin.Context) {
	var req serverOptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commandStart := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && %s start'"`,
		req.LxcId,
		req.User,
		req.Placeholder)

	slog.Info("Command: " + commandStart)

	optStartReturn, err := server.sshClient.NewSession(commandStart)

	if err != nil {
		slog.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	slog.Info(optStartReturn)
	c.IndentedJSON(http.StatusOK, nil)
}
