package env

import (
	"log"
	"os"
	"strings"
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
