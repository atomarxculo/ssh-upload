package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error al cargar la contraseña", err)
	}

	return os.Getenv(key)
}
