package api

import (
	"context"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/sshClient"
	"example/Go-PM-API/util"
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

// API Request body
type ServerRequest struct {
	User string `json:"user" binding:"required"`
}

type Server struct {
	config    util.Config
	ctx       context.Context
	router    *gin.Engine
	pmClient  *proxmoxClient.ProxmoxClient
	sshClient *sshClient.SshClient
}

func NewServer(config util.Config, ctx context.Context, proxmoxClient *proxmoxClient.ProxmoxClient, sshClient *sshClient.SshClient) (*Server, error) {
	server := &Server{
		config:    config,
		ctx:       ctx,
		pmClient:  proxmoxClient,
		sshClient: sshClient,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	defaultLogger := slog.Default()

	// Sets up the slog middleware for GIN
	router.Use(sloggin.New(defaultLogger))

	// Uses the Proxmox API
	router.GET("/containers", server.getContainers)        // Returns info about all containers
	router.GET("/containers/:id", server.getContainerById) // Returns info about a specific container

	// Uses SSH to connect to a container and run the commands
	router.POST("/containers/:id/start", server.postStartServer)     // Start server
	router.POST("/containers/:id/stop", server.postStopServer)       // Stop server
	router.POST("/containers/:id/details", server.postDetailsServer) // Server status -> Online/Offline
	router.POST("/containers/:id/restart", server.postRestartServer) // Restart server

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
