package main

import (
	"context"
	"example/Go-PM-API/api"
	"example/Go-PM-API/proxmoxClient"
	"example/Go-PM-API/util"
	"log"
)

/*
	func pingCtn(ctx context.Context) {
		ctnList, _ := Node.Containers(ctx)
		var totalMem uint64

		for i := 0; i < len(ctnList); i++ {
			ctn := ctnList[i]

			fmt.Println(ctn.Name)
			fmt.Println(ctn.MaxMem / 1048576) // Memória máxima configurada para o container atual - Convertido para MB (de bytes)

			if ctn.Status == "running" { // Caso esteja rodando o container - soma o valor total de memória
				totalMem += ctn.MaxMem / 1048576
			}

			err := Client.Get(ctx, fmt.Sprintf("/nodes/%s/lxc/%d/config", ctn.Node, ctn.VMID), &ctn.ContainerConfig)
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
		privateBytes, err := os.ReadFile(config.SSHKeyFile)
		if err != nil {
			log.Fatal("Failed to load private key (./id_ed25519)")
		}

		signer, err := ssh.ParsePrivateKeyWithPassphrase(privateBytes, []byte(config.SSHKeyPassphrase))
		if err != nil {
			log.Fatal("Failed to parse private key")
		}

		configSSH := &ssh.ClientConfig{
			User:            "root",
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO - Change HostKeyCallback
		}

		client, err := ssh.Dial("tcp", "192.168.18.125:22", configSSH) // TODO - Set a new config for the IP
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
		//	if err := session.Run("/usr/bin/whoami"); err != nil {
		if err := session.Run("lxc-attach -n 101 --uid 1001 /home/untserver/untserver details"); err != nil {
			log.Fatal("Failed to run: " + err.Error())
		}
		fmt.Println(b.String())
	}
*/
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
	println(pmClient)

	// API Server
	server, err := api.NewServer(config, ctx, pmClient)
	println(server)

	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start("localhost:8090")
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
