package main

import (
	"bytes"
	"context"
	"example/Go-PM-API/util"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luthermonson/go-proxmox"
	"golang.org/x/crypto/ssh"
)

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
	var totalMem uint64
	//ctnList[1].Ping()

	for i := 0; i < len(ctnList); i++ {
		ctn := ctnList[i]

		fmt.Println(ctn.Name)
		fmt.Println(ctn.MaxMem / 1048576) // Memória máxima configurada para o container atual - Convertido para MB (de bytes)

		if ctn.Status == "running" { // Caso esteja rodando o container - soma o valor total de memória
			totalMem += ctn.MaxMem / 1048576
		}

		err := Client.Get(Ctx, fmt.Sprintf("/nodes/%s/lxc/%d/config", ctn.Node, ctn.VMID), &ctn.ContainerConfig)
		if err == nil {
			fmt.Println(ctn.ContainerConfig.Hostname)
			fmt.Println(ctn.ContainerConfig.Nets)
			fmt.Println(ctn.ContainerConfig.OnBoot)
		}
		fmt.Println("--------------------------")
	}

	fmt.Printf("\nTotal Memory Used: %dMb | %dGb", totalMem, totalMem/1024)
}

func testSSH(config util.Config) {
	/*

		FAZER FUNCIONAR O SSH COM PUBLIC KEY

	*/

	/*

		SERVER -- https://gist.github.com/jpillora/b480fde82bff51a06238

		privateBytes, err := os.ReadFile(config.SSHKeyFile)
		if err != nil {
			log.Fatal("Failed to load private key (./id_ed25519)")
		}

		private, err := ssh.ParsePrivateKeyWithPassphrase(privateBytes, []byte(config.SSHKeyPassphrase))
		if err != nil {
			log.Fatal("Failed to parse private key")
		}

		configSSH := &ssh.ServerConfig{}
		configSSH.AddHostKey(private)
	*/

	var hostKey ssh.PublicKey
	configSSH := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("Blyanno#1"),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	client, err := ssh.Dial("tcp", "192.168.18.125:22", configSSH)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
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

	//pingCtn()
	testSSH(config)

	// Router
	/*router := gin.Default()
	router.GET("/containers", getContainers)
	router.GET("/containers/:id", getContainerById)
	router.Run("localhost:8090")*/
}
