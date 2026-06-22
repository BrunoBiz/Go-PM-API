package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.GET("/containers", getContainers)
	router.POST("/containers", createContainer)

	router.Run("localhost:8090")
}
