package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// API Request body
type serverStartRequest struct {
	User string `json:"user" binding:"required"`
}

// The output log is formatted as a JSON, which is unmarshaled into this struct
type serverStartOutput struct {
	Option        string `json:"option" binding:"required"`
	Message       string `json:"message" binding:"required"`
	Success       bool   `json:"success" binding:"required"`
	Status        bool   `json:"status" binding:"required"`
	Error         string `json:"error" binding:"required"`
	Command       string `json:"command" binding:"required"`
	CommandResult string `json:"commandresult" binding:"required"`
}

func (server *Server) postStartServer(c *gin.Context) {
	var req serverStartRequest
	var serverOutput serverStartOutput
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

	// Unmarshalls the JSON formatted log into the struct
	err = json.Unmarshal([]byte(optStartReturn[25:]), &serverOutput)

	if err != nil {
		slog.Error("JSON Unmarshal - " + err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err)
	}

	fmt.Println(optStartReturn)

	c.IndentedJSON(http.StatusOK, serverOutput.Message) // TODO - Write a proper return message - Check if already running or started and such -- Use returned message
}
