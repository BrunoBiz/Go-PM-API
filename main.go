package main

import (
	"context"
	"example/Go-PM-API/api"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/sshClient"
	"example/Go-PM-API/util"
	"log"
)

func main() {
	// Loads config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load from config: ", err)
	}

	// Creates context
	ctx := context.Background()

	// Proxmox Client
	pmClient, err := proxmoxClient.NewClient(config, ctx)
	if err != nil {
		log.Fatal("Could not establish proxmox connection", err)
	}

	// SSH Client
	sshClient, err := sshClient.NewSshClient(config)
	if err != nil {
		log.Fatal("Could not establish ssh connection", err)
	}

	// API Server
	server, err := api.NewServer(config, ctx, pmClient, sshClient)
	if err != nil {
		log.Fatal("cannot create server")
	}

	// API server start
	if server.Start("localhost:8090") != nil {
		log.Fatal("cannot start server", err)
	}

	// Close SSH Connection
	sshClient.CloseConnection() // TODO no idea if this is appropriate
}
