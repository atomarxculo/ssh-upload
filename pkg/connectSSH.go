package pkg

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

var (
	client *ssh.Client
	err    error
)

func ConnectSSH(username string, port string, password string, server string, command string, wg *sync.WaitGroup) {
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

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Println("Error al ejecutar " + err.Error())
	}
	fmt.Println("Servidor:", server)
	fmt.Println("Output comando:", b.String())

	wg.Done()

}
