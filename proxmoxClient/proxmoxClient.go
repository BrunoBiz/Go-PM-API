package proxmoxClient

import (
	"context"
	"example/Go-PM-API/util"
	"time"

	"github.com/luthermonson/go-proxmox"
	"golang.org/x/crypto/ssh"
)

type ProxmoxClient struct {
	Client *proxmox.Client
	Node   *proxmox.Node
	config util.Config
	ssh    *ssh.Client
	ctx    context.Context
}

func NewClient(config util.Config, ctx context.Context) (*ProxmoxClient, error) {
	// Proxmox API Client connection
	client := proxmox.NewClient(config.PVEUrl,
		proxmox.WithAPIToken(config.PVEUserRealm+"!"+config.PVETokenID, config.PVEToken),
		proxmox.WithTimeout(30*time.Second), // http.DefaultClient has no timeout
	)

	// Checks if the client is answering
	_, err := client.Version(ctx)
	if err != nil || client == nil {
		return nil, err
	}

	// Main node - This API will load only one node from the proxmox server
	node, err := client.Node(ctx, config.PVENodeName)
	if err != nil {
		return nil, err
	}

	proxmoxClient := &ProxmoxClient{
		Client: client,
		Node:   node,
		config: config,
		ctx:    ctx,
	}

	return proxmoxClient, nil
}
