package api

import (
	"context"
	"example/Go-PM-API/logger"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/sshClient"
	"example/Go-PM-API/util"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

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

	// Sets up HUMA
	humaApi := humagin.New(router, huma.DefaultConfig("Test-api", "1.0.0"))

	// Sets up routing
	// Uses the Proxmox API
	huma.Register(humaApi, huma.Operation{ // Returns info about all containers
		OperationID: "get-containers",
		Method:      http.MethodGet,
		Path:        "/containers",
		Summary:     "Returns info about all containers",
	}, server.getContainers)

	huma.Register(humaApi, huma.Operation{ // Returns info about a specific container
		OperationID: "get-containersById",
		Method:      http.MethodGet,
		Path:        "/containers/:id",
		Summary:     "Returns info about a specific container",
	}, server.getContainerById)

	huma.Register(humaApi, huma.Operation{ // Returns if a specific container is running
		OperationID: "get-containerStatusById",
		Method:      http.MethodGet,
		Path:        "/containers/:id/status",
		Summary:     "Returns if a specific container is running",
	}, server.getContainerStatusById)

	// Uses SSH to connect to a container and run the commands
	huma.Register(humaApi, huma.Operation{ // Start server
		OperationID: "post-startServer",
		Method:      http.MethodPost,
		Path:        "/containers/server/:id/start",
		Summary:     "Starts the container's associated game server",
	}, server.postStartServer)

	huma.Register(humaApi, huma.Operation{ // Server status -> Online/Offline
		OperationID: "post-detailsServer",
		Method:      http.MethodPost,
		Path:        "/containers/server/:id/details",
		Summary:     "Check if the container's associated game server is running",
	}, server.postDetailsServer)

	huma.Register(humaApi, huma.Operation{ // Stop server
		OperationID: "post-stopServer",
		Method:      http.MethodPost,
		Path:        "/containers/server/:id/stop",
		Summary:     "Stops the container's associated game server",
	}, server.postStopServer)

	huma.Register(humaApi, huma.Operation{ // Restart server
		OperationID: "post-restartServer",
		Method:      http.MethodPost,
		Path:        "/containers/server/:id/restart",
		Summary:     "Restarts the container's associated game server",
	}, server.postRestartServer)

	server.router = router
	server.humaAPI = humaApi
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.router)
}

func NewSuccessfulResponse(status bool, message string) ServerResponse {
	response := ServerResponse{}

	response.Body.Successful = true
	response.Body.Status = status
	response.Body.Message = message
	response.Body.ErrorCode = "" // Not used in a SUCCESSFUL response

	return response
}

func NewErrorResponse(message string, errorCode string) ServerResponse { // Will only be considered as an error status codes in the 5xx range
	response := ServerResponse{}

	response.Body.Successful = false
	response.Body.Status = false // Not used in an UNSUCCESSFUL response
	response.Body.Message = message
	response.Body.ErrorCode = errorCode

	return response
}
