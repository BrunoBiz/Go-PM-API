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

func (server *Server) postStopServer(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[postStopServer] - API CALL")
	var req ServerRequest
	var cntID uint64

	// Parameter sent via URL
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	slog.Log(c.Request.Context(), logger.LevelFile, "[postStopServer] - cntID: "+strconv.FormatUint(cntID, 10))

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postStopServer] - ERROR: "+err.Error())
		c.Data(http.StatusBadRequest, "application/json", NewErrorResponse("Bad request.", "BAD_REQUEST"))
		return
	}

	// Prepares the command to start the server
	commandStop := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./gameserver stop'"`,
		cntID,
		req.User)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optStopReturn, err := server.sshClient.NewSession(commandStop)
	slog.Log(c.Request.Context(), logger.LevelFile, "[postStopServer] - commandDetails: "+optStopReturn)

	if err != nil {
		slog.Error("[postStopServer] - SSH New Session: " + err.Error())
		c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
		return
	}

	// Server stopped
	if strings.Contains(optStopReturn, "[  OK  ] Stopping") || strings.Contains(optStopReturn, "MESSAGE: Server stopped") {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postStopServer] - Server stopped successfully")
		c.Data(http.StatusOK, "application/json", NewSuccessfulResponse(true, "Server stopped successfully"))
		return
	}

	// Server is already stopped
	if strings.Contains(optStopReturn, "is already stopped") {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postStopServer] - Server is already stopped")
		c.Data(http.StatusConflict, "application/json", NewSuccessfulResponse(false, "Server is already stopped"))
		return
	}

	// Error
	c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
	slog.Log(c.Request.Context(), logger.LevelFile, "[postStopServer] - OK")
}
