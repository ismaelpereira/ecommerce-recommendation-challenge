package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProjectID        string
	BigtableInstance string
	Port             string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env")
	}

	cfg := &Config{
		ProjectID:        os.Getenv("PROJECT_ID"),
		BigtableInstance: os.Getenv("BIGTABLE_INSTANCE"),
	}
	validate(cfg)

	return cfg
}

func validate(c *Config) {
	if c.ProjectID == "" {
		log.Fatal("PROJECT_ID is required")
	}

	if c.BigtableInstance == "" {
		log.Fatal("BIGTABLE_INSTANCE is required")
	}
}
