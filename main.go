package main

import (
	"context"
	"example/Go-PM-API/api"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/sshClient"
	"example/Go-PM-API/util"
	"fmt"
	"log"
)

func main() {
	// Creates context
	ctx := context.Background()

	// Loads config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load from config: ", err)
	}

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
	fmt.Println(sshClient) // TODO Placeholder
	//sshClient.NewSession("lxc-attach -n 101 --uid 1001 /home/untserver/untserver details")

	// API Server
	server, err := api.NewServer(config, ctx, pmClient)
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
