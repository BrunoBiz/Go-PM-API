package main

import (
	"context"
	"example/Go-PM-API/util"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
var Client *proxmox.Client

func getContainers(c *gin.Context) {

	// Container list in main node
	ctnList, err := Node.Containers(Ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, nil)
	}
	c.IndentedJSON(http.StatusOK, ctnList)
}

func getContainerById(c *gin.Context) {
	var cntID uint64
	cntID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	// Container list in main node
	ctnList, err := Node.Containers(Ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	// Looks for the container

	for i := 0; i < len(ctnList); i++ {
		if uint64(ctnList[i].VMID) == cntID {
			c.IndentedJSON(http.StatusOK, ctnList[i])
		}
	}
	c.IndentedJSON(http.StatusNotFound, nil)

}

func pingCtn() {
	ctnList, _ := Node.Containers(Ctx)

	for i := 0; i < len(ctnList); i++ {
		ctn := ctnList[i]

		fmt.Println(ctn.Name)
		fmt.Println(Client.Get(Ctx, "/nodes/blyanno/lxc/101/config", &ctn.ContainerConfig))
		fmt.Println("--------------------------")

		fmt.Println(ctn.ContainerConfig.Hostname)
		fmt.Println(ctn.ContainerConfig.Nameserver)
		fmt.Println(ctn.ContainerConfig.Nets)

	}
}

func main() {

	// Config loading
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load from config: ", err)
	}

	// ProxMox Connection
	Client = proxmox.NewClient(config.PVEUrl,
		proxmox.WithAPIToken(config.PVEUserRealm+"!"+config.PVETokenID, config.PVEToken),
		proxmox.WithTimeout(30*time.Second), // http.DefaultClient has no timeout
	)

	// Creates context
	Ctx = context.Background()

	// ProxMox Validation
	_, err = Client.Version(Ctx)
	if err != nil || Client == nil {
		log.Fatal("cant connect to proxmox: ", err)
	}

	// Main node - Loads only one
	Node, err = Client.Node(Ctx, config.PVENodeName)
	if err != nil {
		log.Fatal("cant retrieve main node: ", err)
	}

	pingCtn()

	// Router
	/*router := gin.Default()
	router.GET("/containers", getContainers)
	router.GET("/containers/:id", getContainerById)
	router.Run("localhost:8090")*/
}
