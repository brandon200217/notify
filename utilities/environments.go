package utilities

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func GetEnvOrDefault(envName string, defaultValue string) string {
	env := strings.TrimSpace(os.Getenv(envName))
	if env == "" {
		slog.Debug("variable de entorno no seteada, usando default",
			"var", envName,
			"default", defaultValue)
		return defaultValue
	}
	return env
}

func GetEnvOrDefaultInt(envName string, defaultValue int) int {
	env := strings.TrimSpace(os.Getenv(envName))
	if env == "" {
		slog.Debug("variable de entorno no seteada, usando default",
			"var", envName,
			"default", defaultValue)
		return defaultValue
	}
	v, err := strconv.Atoi(env)
	if err != nil {
		slog.Warn("valor numérico inválido en variable de entorno, usando default",
			"var", envName,
			"value", env,
			"default", defaultValue,
			"error", err.Error())
		return defaultValue
	}
	return v
}

func GetEnvOrDefaultBool(envName string, defaultValue bool) bool {
	env := strings.TrimSpace(os.Getenv(envName))
	if env == "" {
		slog.Debug("variable de entorno no seteada, usando default",
			"var", envName,
			"default", defaultValue)
		return defaultValue
	}
	return env == "true"
}

func GetEnvOrDefaultFloat(envName string, defaultValue float64) float64 {
	env := strings.TrimSpace(os.Getenv(envName))
	if env == "" {
		log.Printf("envName: %s - Using default: %v", envName, defaultValue)
		return defaultValue
	}
	v, err := strconv.ParseFloat(env, 64)
	if err != nil {
		log.Printf("envName: %s - value: %s - Invalid float, using default: %v", envName, env, defaultValue)
		return defaultValue
	}
	return v
}
