package api

import (
	"example/Go-PM-API/logger"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/acarl005/stripansi"
	"github.com/gin-gonic/gin"
)

func (server *Server) postDetailsServer(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[postDetailsServer] - API CALL")

	var req ServerRequest
	var cntID uint64

	// Parameter sent via URL
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	slog.Log(c.Request.Context(), logger.LevelFile, "[postDetailsServer] - cntID: "+strconv.FormatUint(cntID, 10))

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postDetailsServer] - ERROR: "+err.Error())
		c.Data(http.StatusBadRequest, "application/json", NewErrorResponse("Bad request.", "BAD_REQUEST"))
		return
	}

	// Prepares the command to check server details
	commandDetails := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./gameserver details'"`,
		cntID,
		req.User)
	slog.Log(c.Request.Context(), logger.LevelFile, "[postDetailsServer] - commandDetails: "+commandDetails)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optDetailsReturn, err := server.sshClient.NewSession(commandDetails)

	if err != nil {
		slog.Error("[postDetailsServer] - SSH New Session: " + err.Error())
		c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
		return
	}

	// Server ONLINE
	if regexp.MustCompile(`(?mi)(status:)\s+(started)`).MatchString(stripansi.Strip(optDetailsReturn)) {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postDetailsServer] - Server running")
		c.Data(http.StatusOK, "application/json", NewSuccessfulResponse(true, "Server running"))
		return
	}

	// Server OFFLINE
	if regexp.MustCompile(`(?mi)(status:)\s+(stopped)`).MatchString(stripansi.Strip(optDetailsReturn)) {
		slog.Log(c.Request.Context(), logger.LevelFile, "[postDetailsServer] - Server stopped")
		c.Data(http.StatusOK, "application/json", NewSuccessfulResponse(true, "Server stopped"))
		return
	}

	// Error
	c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
	slog.Log(c.Request.Context(), logger.LevelFile, "[postDetailsServer] - OK")
}
