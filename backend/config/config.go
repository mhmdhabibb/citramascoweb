package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {

	err := godotenv.Load()
	if err != nil {
		// This will print a warning in your terminal if it can't find the file
		println("Warning: .env file not found!")
	}
	return os.Getenv(key)

}
