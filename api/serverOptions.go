package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type serverOptRequest struct {
	LxcId       int    `json:"lxcId" binding:"required"`
	UserId      int    `json:"uid" binding:"required"`
	Placeholder string `json:"Placeholder" binding:"required"` // TODO
}

func (server *Server) postStartServer(c *gin.Context) {
	var req serverOptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO - SSH HERE - LXC ATTACH -> START SERVER

}
