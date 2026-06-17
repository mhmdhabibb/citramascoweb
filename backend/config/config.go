package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("ENV not found!")
	}

	return os.Getenv(key)

}
