package api

import (
	"context"
	"encoding/json"
	"example/Go-PM-API/logger"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/sshClient"
	"example/Go-PM-API/util"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/luthermonson/go-proxmox"
	sloggin "github.com/samber/slog-gin"
)

// API Request Body
type ServerRequest struct {
	User string `json:"user" binding:"required"`
}

// API response body
type serverResponse struct {
	Successful bool   `json:"successful" binding:"required"` // If the request was successful or not - False when an error occurs
	Status     bool   `json:"status" binding:"required"`     // If the request could achieve what was intended / If the request is based on a boolean response E.g. If the server is running - True / If the container was not found - False
	Message    string `json:"message" binding:"required"`
	ErrorCode  string `json:"errorcode" binding:"required"`
}

// Container info from the proxmox package, formatted to be used with HUMA
type ContainerInfo struct {
	CPUs    int                    `json:"cpus"`
	MaxDisk uint64                 `json:"maxdisk"`
	MaxMem  uint64                 `json:"maxmem"`
	MaxSwap uint64                 `json:"maxswap"`
	Name    string                 `json:"name"`
	Node    string                 `json:"node"`
	Status  string                 `json:"status"`
	Tags    string                 `json:"tags"`
	Uptime  uint64                 `json:"uptime"`
	VMID    proxmox.StringOrUint64 `json:"vmid"`
}

// An array of containers
type ContainerOutput struct {
	Body struct {
		Containers []ContainerInfo
	}
}

type Server struct {
	config    util.Config
	ctx       context.Context
	router    *gin.Engine
	pmClient  *proxmoxClient.ProxmoxClient
	sshClient *sshClient.SshClient
	humaAPI   huma.API
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

	// Sets up HUMA
	humaApi := humagin.New(router, huma.DefaultConfig("Test-api", "1.0.0"))

	// Sets up routing
	// Uses the Proxmox API
	huma.Register(humaApi, huma.Operation{
		OperationID: "get-containers",
		Method:      http.MethodGet,
		Path:        "/containers",
		Summary:     "Get a list of all containers",
	}, server.getContainers)

	huma.Register(humaApi, huma.Operation{
		OperationID: "get-containersById",
		Method:      http.MethodGet,
		Path:        "/containers/:id",
		Summary:     "Get a contaiber by ID",
	}, server.getContainerById)

	/*
		router.GET("/containers", server.getContainers)                     // Returns info about all containers
		router.GET("/containers/:id", server.getContainerById)              // Returns info about a specific container
		router.GET("/containers/:id/status", server.getContainerStatusById) // Returns if a specific container is running
	*/

	// Uses SSH to connect to a container and run the commands
	/*
		router.POST("/containers/server/:id/start", server.postStartServer)     // Start server
		router.POST("/containers/server/:id/stop", server.postStopServer)       // Stop server
		router.POST("/containers/server/:id/details", server.postDetailsServer) // Server status -> Online/Offline
		router.POST("/containers/server/:id/restart", server.postRestartServer) // Restart server
	*/
	server.router = router
	server.humaAPI = humaApi
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.router)
	//return server.router.Run(address)
}

func NewSuccessfulResponse(status bool, message string) []byte {
	response := serverResponse{
		Successful: true,
		Status:     status,
		Message:    message,
		ErrorCode:  "", // Not used in a SUCCESSFUL response
	}

	marshaledResponse, _ := json.Marshal(response)

	return marshaledResponse
}

func NewErrorResponse(message string, errorCode string) []byte { // Will only be considered as an error status codes in the 5xx range
	response := serverResponse{
		Successful: false,
		Status:     false, // Not used in an UNSUCCESSFUL response
		Message:    message,
		ErrorCode:  errorCode,
	}

	marshaledResponse, _ := json.Marshal(response)

	return marshaledResponse
}
