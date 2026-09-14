package api

import (
	"encoding/json"
	"example/Go-PM-API/logger"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) getContainers(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainers] - API CALL")

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainers] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[getContainers] - ERROR: "+err.Error())
		requestErrorMsg, _ := json.Marshal(ServerError{Message: "An error occurred while processing the request.", Code: "INTERNAL_SERVER_ERROR"})
		c.IndentedJSON(http.StatusInternalServerError, string(requestErrorMsg))
		return
	}

	c.IndentedJSON(http.StatusOK, ctnList)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainers] - OK")
}

func (server *Server) getContainerById(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - API CALL")

	var cntID uint64
	var serverResponse = ServerResponse{Status: false, Message: "No container found"} // Default response - no container found

	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - cntID:"+strconv.FormatUint(cntID, 10))

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - ERROR: "+err.Error())
		requestErrorMsg, _ := json.Marshal(ServerError{Message: "An error occurred while processing the request.", Code: "INTERNAL_SERVER_ERROR"})
		c.IndentedJSON(http.StatusInternalServerError, string(requestErrorMsg))
		return
	}

	// Looks for the container
	for i := 0; i < len(ctnList); i++ {
		if uint64(ctnList[i].VMID) == cntID {
			slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - Container found")
			c.IndentedJSON(http.StatusOK, ctnList[i])
			slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - OK")
			return
		}
	}

	c.IndentedJSON(http.StatusOK, serverResponse)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - OK")
}

func (server *Server) getContainerStatusById(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - API CALL")

	var cntID uint64
	var serverResponse = ServerResponse{Status: false, Message: "No container found"} // Default response - no container found

	cntID, _ = strconv.ParseUint(c.Param("id"), 10, 64)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - cntID:"+strconv.FormatUint(cntID, 10))

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - ERROR: "+err.Error())
		requestErrorMsg, _ := json.Marshal(ServerError{Message: "An error occurred while processing the request.", Code: "INTERNAL_SERVER_ERROR"})
		c.IndentedJSON(http.StatusInternalServerError, string(requestErrorMsg))
		return
	}

	// Looks for the container
	for i := 0; i < len(ctnList); i++ {
		if uint64(ctnList[i].VMID) == cntID {
			slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - Container found")
			serverResponse = ServerResponse{Status: true, Message: ctnList[i].Status}
		}
	}

	serverResponseJSON, _ := json.Marshal(serverResponse)
	c.IndentedJSON(http.StatusOK, string(serverResponseJSON))
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - OK")
}
