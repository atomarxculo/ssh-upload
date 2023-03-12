package main

import (
	"flag"
	"fmt"
	"os"
	"ssh-upload/pkg"
	"strings"
	"sync"
)

func main() {
	var server, username, port, password, command, localfile, remotefile string
	var wg = &sync.WaitGroup{}

	flag.StringVar(&username, "user", "", "Indicar el nombre del usuario con el que queramos conectarnos.")
	flag.StringVar(&port, "port", "22", "Puerto con el que queramos conectarnos.")
	flag.StringVar(&server, "server", "", "Indicar el servidor al que queramos conectarnos, si se deja en blanco, lee el fichero 'servers.txt' que se encuentre del mismo directorio.")
	flag.StringVar(&password, "pass", "", "Contraseña, por defecto lee la variable PASSWORD en el .env que se encuentre del mismo directorio.")
	flag.StringVar(&command, "command", "hostname", "Comando que va a ejecutar en el servidor.")
	flag.StringVar(&localfile, "localfile", "", "Fichero local a subir")
	flag.StringVar(&remotefile, "remotefile", "/var/tmp/", "Destino del fichero")
	flag.Parse()

	env := pkg.GetEnvVariable("PASSWORD")
	if !pkg.FlagPassed("server") && pkg.FlagPassed("port") {
		fmt.Println("Si indicas el puerto, también tienes que indicar el servidor.")
		os.Exit(0)
	} else if !pkg.FlagPassed("server") && !pkg.FlagPassed("user") {
		fmt.Println("Tienes que indicar un usuario para poder conectarte.")
		os.Exit(0)
	} else if pkg.FlagPassed("remotefile") && !pkg.FlagPassed("localfile") {
		fmt.Println("Tienes que elegir un fichero local que subir")
		os.Exit(0)
	}

	if pkg.FlagPassed("server") {
		wg.Add(1)
		go pkg.ConnectSSH(username, port, env, server, command, localfile, remotefile, wg)
		wg.Wait()
	} else {
		if strings.Contains(server, ":") {
			port = ""
		}
		servertext, err := pkg.ReadServer("servers.txt")
		if err != nil {
			fmt.Println(err)
		}
		for _, serverline := range servertext {
			wg.Add(1)
			go pkg.ConnectSSH(username, port, env, serverline, command, localfile, remotefile, wg)
		}
		wg.Wait()
	}
}
