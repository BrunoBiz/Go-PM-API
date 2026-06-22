package main

import (
	"context"
	"example/Go-PM-API/util"
	"fmt"
	"log"
	"net/http"

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

func getContainers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, containers)
}

func createContainer(c *gin.Context) {
	var newContainer container

	if err := c.BindJSON(&newContainer); err != nil {
		return
	}

	containers = append(containers, newContainer)
	c.IndentedJSON(http.StatusCreated, newContainer)
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
		//proxmox.WithTimeout(30*time.Second), // http.DefaultClient has no timeout
	)

	// ProxMox Validation
	ctx := context.Background()
	version, err := client.Version(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(version.Release) // 6.3

	fmt.Println("teste")

	node, err := client.Node(ctx, "blyanno")
	fmt.Println(node.Name)
	fmt.Println(node.Kversion)
	fmt.Println(node.LoadAvg)
	fmt.Println(node.CPU)
	fmt.Println(node.RootFS)
	fmt.Println(node.PVEVersion)
	fmt.Println(node.CPUInfo)
	fmt.Println(node.Swap)
	fmt.Println(node.Idle)
	fmt.Println(node.Memory)
	fmt.Println(node.Ksm)
	fmt.Println(node.Uptime)
	fmt.Println(node.Wait)
	fmt.Println(err)

	// Router
	/*
		router := gin.Default()
		router.GET("/containers", getContainers)
		router.POST("/containers", createContainer)
	*/
	//router.Run("localhost:8090")
}
