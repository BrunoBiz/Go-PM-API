package main

import (
	"context"
	"example/Go-PM-API/util"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luthermonson/go-proxmox"
)

type container struct {
	ID        int    `json:"id"`
	ProxMoxID int    `json:"proxmoxid"`
	Name      string `json:"name"`
	Status    bool   `json:"status"`
}

var containers = []container{
	{ID: 1, ProxMoxID: 100, Name: "MBC-MineServer", Status: true},
	{ID: 2, ProxMoxID: 101, Name: "UnturnedServer", Status: false},
	{ID: 3, ProxMoxID: 103, Name: "Kanboard", Status: true},
}

var Node *proxmox.Node
var Ctx context.Context
var ctnList proxmox.Containers

func getContainers(c *gin.Context) {

	// Container list in main node
	ctnList, err := Node.Containers(Ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, nil)
	}
	c.IndentedJSON(http.StatusOK, ctnList)
}

func main() {

	// Config loading
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load from config: ", err)
	}

	// ProxMox Connection
	client := proxmox.NewClient(config.PVEUrl,
		proxmox.WithAPIToken(config.PVEUserRealm+"!"+config.PVETokenID, config.PVEToken),
		proxmox.WithTimeout(30*time.Second), // http.DefaultClient has no timeout
	)

	// Creates context
	Ctx = context.Background()

	// ProxMox Validation
	_, err = client.Version(Ctx)
	if err != nil || client == nil {
		log.Fatal("cant connect to proxmox: ", err)
	}

	// Main node - Loads only one
	Node, err = client.Node(Ctx, config.PVENodeName)
	if err != nil {
		log.Fatal("cant retrieve main node: ", err)
	}

	// Router
	router := gin.Default()
	router.GET("/containers", getContainers)
	router.Run("localhost:8090")
}
