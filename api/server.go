package api

import (
	"context"
	"example/Go-PM-API/util"

	"github.com/gin-gonic/gin"
	"github.com/luthermonson/go-proxmox"
)

type Server struct {
	config util.Config
	client *proxmox.Client
	ctx    context.Context
	node   *proxmox.Node
	router *gin.Engine
}

func NewServer(config util.Config, client *proxmox.Client, ctx context.Context, node *proxmox.Node) (*Server, error) {
	server := &Server{
		config: config,
		client: client,
		ctx:    ctx,
		node:   node,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/containers", server.getContainers)
	router.GET("/containers/:id", server.getContainerById)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
