package sshClient

import (
	"bytes"
	"example/Go-PM-API/util"
	"os"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	sshClient *ssh.Client
	config    util.Config
}

func NewSshClient(config util.Config) (*SshClient, error) {
	// Reads private key
	privateBytes, err := os.ReadFile(config.SSHKeyFile)
	if err != nil {
		return nil, err
		// log.Fatal("Failed to load private key (./id_ed25519)")
	}

	// Creates signer
	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateBytes, []byte(config.SSHKeyPassphrase))
	if err != nil {
		return nil, err
		// log.Fatal("Failed to parse private key")
	}

	// SSH client config
	configSSH := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO - Change HostKeyCallback
	}

	client, err := ssh.Dial("tcp", config.SSHPveIP+":"+config.SSHPvePort, configSSH)
	if err != nil {
		return nil, err
		//log.Fatal("Failed to dial: ", err)
	}

	sshClient := &SshClient{
		sshClient: client,
		config:    config,
	}

	// defer client.Close() 	// TODO - Close connection somewhere
	return sshClient, nil
}

func (ssh *SshClient) NewSession(context string) (string, error) {
	session, err := ssh.sshClient.NewSession()
	if err != nil {
		return "", err
		// log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	//"lxc-attach -n 101 --uid 1001 /home/untserver/untserver details"

	if err := session.Run(context); err != nil {
		return "", err
		//log.Fatal("Failed to run: " + err.Error())
	}

	return b.String(), nil
}

func (ssh *SshClient) CloseConnection() {
	ssh.sshClient.Close()
}
