package pkg

import (
	"bufio"
	"log"
	"os"
)

func ReadServer(file string) (fileLines []string, err error) {
	readFile, err := os.Open(file)
	if err != nil {
		log.Fatal("Error al leer el fichero", err)
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	err = fileScanner.Err()

	return
}
