package utilities

import (
	"log"
	"os"
	"strconv"
)

func GetEnvOrDefaultInt(envName string, defaultValue int) int {
	if env := os.Getenv(envName); env == "" {

		log.Printf("envName: %s - value: %s - Using default: %v", envName, env, defaultValue)
		return defaultValue
	} else {
		v, err := strconv.Atoi(env)
		if err != nil {
			log.Printf("envName: %s - value: %s - Invalid numeric value, using default: %v", envName, env, defaultValue)
			return defaultValue
		}
		return v
	}
}

func GetEnvOrDefault(envName string, defaultValue string) string {
	if env := os.Getenv(envName); env == "" {
		log.Printf("envName: %s - value: %s - Using default: %s", envName, env, defaultValue)
		return defaultValue
	} else {
		return env
	}
}

func GetEnvOrDefaultBool(envName string, defaultValue bool) bool {
	if env := os.Getenv(envName); env == "" {
		log.Printf("envName: %s - value: %s - Using default: %v", envName, env, defaultValue)
		return defaultValue
	} else {
		return env == "true"
	}
}
