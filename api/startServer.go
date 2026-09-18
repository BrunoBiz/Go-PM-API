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

func (server *Server) postStartServer(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[postStartServer] - API CALL")
	var req ServerRequest
	var cntID uint64

	// Parameter sent via URL
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	slog.Log(c.Request.Context(), logger.LevelFile, "[postStartServer] - cntID: "+strconv.FormatUint(cntID, 10))

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postStartServer] - ERROR: "+err.Error())
		c.Data(http.StatusBadRequest, "application/json", NewErrorResponse("Bad request.", "BAD_REQUEST"))
		return
	}

	// Prepares the command to start the server
	commandStart := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./gameserver start'"`,
		cntID,
		req.User)
	slog.Log(c.Request.Context(), logger.LevelFile, "[postStartServer] - commandDetails: "+commandStart)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optStartReturn, err := server.sshClient.NewSession(commandStart)

	if err != nil {
		slog.Error("[postStartServer] - SSH New Session: " + err.Error())
		c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
		return
	}

	// Started
	if strings.Contains(optStartReturn, "[  OK  ] Starting") || strings.Contains(optStartReturn, "MESSAGE: Server started") {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postStartServer] - Server started successfully")
		c.Data(http.StatusOK, "application/json", NewSuccessfulResponse(true, "Server started successfully."))
		return
	}

	// Already running
	if strings.Contains(optStartReturn, "is already running") {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postStartServer] - Server is already running")
		c.Data(http.StatusOK, "application/json", NewSuccessfulResponse(false, "Server is already running."))
		return
	}

	// Error
	c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
	slog.Log(c.Request.Context(), logger.LevelFile, "[postStartServer] - OK")
}
