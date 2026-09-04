package main

import (
	"context"
	"example/Go-PM-API/api"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/sshClient"
	"example/Go-PM-API/util"
	"log"
	"log/slog"
)

func main() {
	// Loads config
	slog.Info("[Starting API] - Loading configuration file...")
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("[Starting API] - Error - Could not load from config: ", err)
	}

	// Creates context
	slog.Info("[Starting API] - Creating default context...")
	ctx := context.Background()

	// Proxmox Client
	slog.Info("[Starting API] - Establishing Proxmox connection...")
	pmClient, err := proxmoxClient.NewClient(config, ctx)
	if err != nil {
		log.Fatal("[Starting API] - Error - Could not establish Proxmox connection", err)
	}

	// SSH Client
	slog.Info("[Starting API] - Establishing SSH connection...")
	sshClient, err := sshClient.NewSshClient(config)
	if err != nil {
		log.Fatal("[Starting API] - Error - Could not establish ssh connection ", err)
	}

	// Close SSH Connection
	defer sshClient.CloseConnection()

	// API Server
	slog.Info("[Starting API] - Creating API server...")
	server, err := api.NewServer(config, ctx, pmClient, sshClient)
	if err != nil {
		log.Fatal("[Starting API] - Error - Could not create server")
	}

	// API server start
	slog.Info("[Starting API] - Starting API server...")
	if server.Start("0.0.0.0:8090") != nil { // TODO - Listen on public and local
		log.Fatal("[Starting API] - Error - Could not start server", err)
	}
}
