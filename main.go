package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"ssh-upload/pkg"

	"golang.org/x/crypto/ssh"
)

func main() {
	var server, username, port, password, command, serverline string

	flag.StringVar(&username, "user", "", "Indicar el nombre del usuario con el que queramos conectarnos.")
	flag.StringVar(&port, "port", "22", "Puerto con el que queramos conectarnos.")
	flag.StringVar(&server, "server", "", "Indicar el servidor al que queramos conectarnos.")
	flag.StringVar(&password, "pass", "", "Contraseña, por defecto lee la variable PASSWORD en el .env que se encuentre del mismo directorio.")
	flag.StringVar(&command, "command", "hostname", "Comando que va a ejecutar en el servidor")
	flag.Parse()

	env := pkg.GetEnvVariable("PASSWORD")

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(env),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var client *ssh.Client
	var err error

	servertext, err := pkg.ReadServer("servers.txt")
	for _, line := range servertext {
		fmt.Println(line)
	}
	fmt.Println(servertext)

	// Meter toda la conexión en una función para que se ejecute en el for

	if server != "" {
		client, err = ssh.Dial("tcp", server+":"+port, config)
	} else {
		client, err = ssh.Dial("tcp", serverline, config)
	}

	if err != nil {
		log.Fatal("Fallo al abrir conexión", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Error al crear la sesión", err)
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatal("Error al ejecutar" + err.Error())
	}
	fmt.Println(b.String())

}
