package env

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	ENV_DB_URI = "DB_URI"
)

var envVars map[string]string

func init() {
	log.Printf("[DEBUG] Loading Environment Variables")

	// Load Environment Variables
	envVars = map[string]string{}
	envs := os.Environ()
	for _, e := range envs {
		split := strings.Split(e, "=")
		envVars[split[0]] = split[1]
	}
	log.Printf("[DEBUG] Environment variables loaded")

	log.Printf("[DEBUG] Environment variables loading from .env")
	env, err := godotenv.Read()
	if err != nil {
		log.Printf("[WARNING] Error loading .env file: %s", err.Error())
	} else {
		for k, v := range env {
			envVars[k] = v
		}
	}
}

func Get(key string) string {
	return envVars[key]
}

func GetDBURI() string {
	uri, ok := envVars[ENV_DB_URI]
	if !ok {
		log.Fatalf("Failed to retrieve %s env", ENV_DB_URI)
	}

	return uri
}
