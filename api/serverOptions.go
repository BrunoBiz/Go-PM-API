package api

import (
	"log/slog"
	"net/http"
	"strconv"

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

	teste, err := server.sshClient.NewSession("lxc-attach -n 106 --uid 0 /opt/gameserver stop")

	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info(teste)
	slog.Info(req.Placeholder)
	slog.Info(strconv.Itoa(req.UserId))
	slog.Info(strconv.Itoa(req.LxcId))

	// TODO - SSH HERE - LXC ATTACH -> START SERVER

}
