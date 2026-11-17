package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	envOnce sync.Once
)

func loadEnv() {
	envOnce.Do(func() {
		// Attempt to load .env from current directory, backend directory, or parent directory
		paths := []string{".env", "backend/.env", "../.env"}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				if err := godotenv.Load(p); err == nil {
					return
				}
			}
		}
		println("Warning: .env file not found!")
	})
}

func Config(key string) string {
	loadEnv()
	return os.Getenv(key)
}
