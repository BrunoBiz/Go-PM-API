package api

import (
	"context"
	"example/Go-PM-API/logger"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/sshClient"
	"example/Go-PM-API/util"
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

// API Request Body
type ServerRequest struct {
	User string `json:"user" binding:"required"`
}

// API Error
type ServerError struct {
	Message string `json:"message" binding:"required"`
	Code    string `json:"code" binding:"required"`
}

// API Response
type ServerResponse struct {
	Status  bool   `json:"status" binding:"required"`
	Message string `json:"message" binding:"required"`
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

	configSlogGin := sloggin.Config{
		DefaultLevel: logger.LevelGIN,
		//ClientErrorLevel: slog.LevelWarn,
		//ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    true,
		WithRequestHeader:  true,
		WithResponseBody:   true,
		WithResponseHeader: true,
		WithSpanID:         false,
		WithTraceID:        false,
		WithClientIP:       false,

		WithCustomMessage: func(c *gin.Context) string {
			return "GIN-Log"
		},
	}

	// Sets up the slog middleware for GIN
	router.Use(sloggin.NewWithConfig(defaultLogger, configSlogGin))

	// Uses the Proxmox API
	router.GET("/containers", server.getContainers)                     // Returns info about all containers
	router.GET("/containers/:id", server.getContainerById)              // Returns info about a specific container
	router.GET("/containers/:id/status", server.getContainerStatusById) // Returns if a specific container is running

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
