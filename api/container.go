package api

import (
	"context"
	"example/Go-PM-API/logger"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*func (server *Server) Asd(ctx context.Context, input *struct{}) (*ContainerOutput, error) {
	resp := &ContainerOutput{}
	return resp, nil
}*/

func (server *Server) getContainers(c context.Context, input *struct{}) (*ContainerOutput, error) {
	slog.Log(c, logger.LevelFile, "[getContainers] - API CALL")

	var containerInfo ContainerInfo
	var containerOutput ContainerOutput

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c, logger.LevelFile, "[getContainers] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c, logger.LevelFile, "[getContainers] - ERROR: "+err.Error())
		//c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
		return nil, err
	}

	for _, ctnRange := range ctnList {
		containerInfo.CPUs = ctnRange.CPUs
		containerInfo.MaxDisk = ctnRange.MaxDisk
		containerInfo.MaxMem = ctnRange.MaxMem
		containerInfo.MaxSwap = ctnRange.MaxSwap
		containerInfo.Name = ctnRange.Name
		containerInfo.Node = ctnRange.Node
		containerInfo.Status = ctnRange.Status
		containerInfo.Tags = ctnRange.Tags
		containerInfo.Uptime = ctnRange.Uptime
		containerInfo.VMID = ctnRange.VMID

		containerOutput.Body.Containers = append(containerOutput.Body.Containers, containerInfo)
	}

	//c.IndentedJSON(http.StatusOK, ctnList)
	slog.Log(c, logger.LevelFile, "[getContainers] - OK")
	return &containerOutput, err
}

func (server *Server) getContainerById(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - API CALL")

	var cntID uint64
	var serverResponse = NewSuccessfulResponse(false, "No container found") // Default response - no container found

	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - cntID:"+strconv.FormatUint(cntID, 10))

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - ERROR: "+err.Error())
		c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
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

	c.Data(http.StatusOK, "application/json", serverResponse)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerById] - OK")
}

func (server *Server) getContainerStatusById(c *gin.Context) {
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - API CALL")

	var cntID uint64
	var serverResponse = NewSuccessfulResponse(false, "No container found") // Default response - no container found

	cntID, _ = strconv.ParseUint(c.Param("id"), 10, 64)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - cntID:"+strconv.FormatUint(cntID, 10))

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - ERROR: "+err.Error())
		c.Data(http.StatusInternalServerError, "application/json", NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR"))
		return
	}

	// Looks for the container
	for i := 0; i < len(ctnList); i++ {
		if uint64(ctnList[i].VMID) == cntID {
			slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - Container found")
			serverResponse = NewSuccessfulResponse(true, ctnList[i].Status)
		}
	}

	c.Data(http.StatusOK, "application/json", serverResponse)
	slog.Log(c.Request.Context(), logger.LevelFile, "[getContainerStatusById] - OK")
}
