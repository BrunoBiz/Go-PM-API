package api

import (
	"context"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config   util.Config
	ctx      context.Context
	router   *gin.Engine
	pmClient *proxmoxClient.ProxmoxClient
}

func NewServer(config util.Config, ctx context.Context, proxmoxClient *proxmoxClient.ProxmoxClient) (*Server, error) {
	server := &Server{
		config:   config,
		ctx:      ctx,
		pmClient: proxmoxClient,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/containers", server.getContainers)
	router.GET("/containers/:id", server.getContainerById)
	router.POST("/server/start", server.postStartServer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
