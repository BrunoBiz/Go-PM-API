package api

import (
	"example/Go-PM-API/logger"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *Server) postRestartServer(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[postRestartServer] - API CALL")
	var req ServerRequest
	var cntID uint64

	// Parameter sent via URL
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	slog.Log(c.Request.Context(), logger.LevelFile, "[postRestartServer] - cntID: "+strconv.FormatUint(cntID, 10))

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postRestartServer] - ERROR: "+err.Error())
		c.Data(http.StatusBadRequest, "application/json", NewErrorResponse("Bad request.", "BAD_REQUEST"))
		return
	}

	// Prepares the command to restart the server
	commandRestart := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./gameserver restart'"`,
		cntID,
		req.User)
	slog.Log(c.Request.Context(), logger.LevelFile, "[postRestartServer] - commandDetails: "+commandRestart)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optRestartReturn, err := server.sshClient.NewSession(commandRestart)

	if err != nil {
		slog.Error("SSH New Session: " + err.Error())
		c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
		return
	}

	// Started
	if strings.Contains(optRestartReturn, "[  OK  ] Starting") || strings.Contains(optRestartReturn, "MESSAGE: Server started") {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postRestartServer] - Server started successfully")
		c.Data(http.StatusOK, "application/json", NewSuccessfulResponse(true, "Server started successfully."))
		return
	}

	// Already running
	if strings.Contains(optRestartReturn, "is already running") {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postRestartServer] - Server is already running")
		c.Data(http.StatusConflict, "application/json", NewSuccessfulResponse(false, "Server is already running."))
		return
	}

	// Error
	c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
	slog.Log(c.Request.Context(), logger.LevelFile, "[postRestartServer] - OK")
}
