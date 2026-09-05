package main

import (
	"context"
	"example/Go-PM-API/api"
	"example/Go-PM-API/logger"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/sshClient"
	"example/Go-PM-API/util"
	"log/slog"
)

func main() {
	// Defines default logger with slog
	slog.Info("[Starting API] - Setting up default logger...")
	err := logger.LoadLogger()
	if err != nil {
		slog.Error("[Starting API] - Could not set a default logger: " + err.Error())
		return
	}

	// Loads config
	slog.Info("[Starting API] - Loading configuration file...")
	config, err := util.LoadConfig(".")
	if err != nil {
		slog.Error("[Starting API] - Could not load from config: " + err.Error())
		return
	}

	// Creates context
	slog.Info("[Starting API] - Creating default context...")
	ctx := context.Background()

	// Proxmox Client
	slog.Info("[Starting API] - Establishing Proxmox connection...")
	pmClient, err := proxmoxClient.NewClient(config, ctx)
	if err != nil {
		slog.Error("[Starting API] - Could not establish Proxmox connection: " + err.Error())
		return
	}

	// SSH Client
	slog.Info("[Starting API] - Establishing SSH connection...")
	sshClient, err := sshClient.NewSshClient(config)
	if err != nil {
		slog.Error("[Starting API] - Could not establish ssh connection: " + err.Error())
		return
	}

	// Close SSH Connection
	defer sshClient.CloseConnection()

	// API Server
	slog.Info("[Starting API] - Creating API server...")
	server, err := api.NewServer(config, ctx, pmClient, sshClient)
	if err != nil {
		slog.Error("[Starting API] - Could not create server: " + err.Error())
		return
	}

	// API server start
	slog.Info("[Starting API] - Starting API server...")
	err = server.Start("0.0.0.0:8090")
	if err != nil {
		slog.Error("[Starting API] - Could not start server: " + err.Error())
		return
	}
}
