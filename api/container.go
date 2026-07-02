package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) getContainers(c *gin.Context) {
	// Container list in main node
	ctnList, err := server.node.Containers(server.ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, nil)
	}
	c.IndentedJSON(http.StatusOK, ctnList)
}

func (server *Server) getContainerById(c *gin.Context) {
	var cntID uint64
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	// Container list in main node
	ctnList, err := server.node.Containers(server.ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	// Looks for the container

	for i := 0; i < len(ctnList); i++ {
		if uint64(ctnList[i].VMID) == cntID {
			c.IndentedJSON(http.StatusOK, ctnList[i])
		}
	}
	c.IndentedJSON(http.StatusNotFound, nil)

}
