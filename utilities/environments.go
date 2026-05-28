package utilities

import (
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
