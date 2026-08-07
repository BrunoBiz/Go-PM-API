package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type serverStartRequest struct {
	User string `json:"user" binding:"required"`
}

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

	if err != nil {
		slog.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	// TESTING

	//slog.Info("SSH combined output: \n" + optStartReturn)
	//slog.Info(optStartReturn)
	//fmt.Println(optStartReturn)

	var serverOutput serverStartOutput
	err = json.Unmarshal([]byte(optStartReturn[25:]), &serverOutput)

	if err != nil {
		slog.Error("Error parsing JSON", "error", err)
	}
	slog.Info(serverOutput.Message)
	slog.Info(serverOutput.Option)
	slog.Info(serverOutput.Command)

	// TESTING

	c.IndentedJSON(http.StatusOK, serverOutput.Message) // TODO - Write a proper return message - Check if already running or started and such -- Use returned message
}
