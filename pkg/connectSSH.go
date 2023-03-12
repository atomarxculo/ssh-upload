package pkg

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

var (
	client *ssh.Client
	err    error
)

func ConnectSSH(username string, port string, password string, server string, command string, localfile string, remotepath string, wg *sync.WaitGroup) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if !strings.Contains(server, ":") {
		client, err = ssh.Dial("tcp", server+":"+port, config)
	} else {
		client, err = ssh.Dial("tcp", server, config)
	}

	if err != nil {
		log.Fatal("Fallo al abrir conexión ", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Error al crear la sesión ", err)
	}

	defer session.Close()

	if command != "" {
		commandSSH(session, command, server)
	}
	if localfile != "" {
		uploadFileSCP(config, server, localfile, remotepath)
	}

	wg.Done()

}

func commandSSH(session *ssh.Session, command string, server string) {
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Println("Error al ejecutar " + err.Error())
	}
	fmt.Println("Servidor:", server)
	fmt.Println("Output comando:", b.String())
}

func uploadFileSCP(config *ssh.ClientConfig, server string, localfile string, remotepath string) {
	fullpath := remotepath + localfile
	client := scp.NewClient(server, config)
	err := client.Connect()
	if err != nil {
		fmt.Println()
	}
	f, _ := os.Open(localfile)
	defer client.Close()
	defer f.Close()
	err = client.CopyFromFile(context.Background(), *f, fullpath, "0644")
	if err != nil {
		fmt.Println(err)
	}
}
