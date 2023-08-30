package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	envPgHost     = "POSTGRES_HOST"
	envPgPort     = "POSTGRES_PORT"
	envPgDB       = "POSTGRES_DB"
	envPgUser     = "POSTGRES_USER"
	envPgPassword = "POSTGRES_PASSWORD"
	envPgSslMode  = "POSTGRES_SSL_MODE"
)

const (
	defaultPgPort    = 5432
	defaultPgDb      = "jump"
	defaultPgSslMode = "disable"
)

type PgConfig struct {
	PgHost     string
	PgPort     int64
	PgDbName   string
	PgUser     string
	PgPassword string
	PgSslMode  string // https://www.postgresql.org/docs/current/libpq-ssl.html#LIBPQ-SSL-SSLMODE-STATEMENTS
}

const (
	envAppHost = "WS_HOST"
	envAppPort = "WS_PORT"
	envAppMode = "WS_MODE"
)

const (
	defaultAppHost = "localhost"
	defaultAppPort = 8080
	defaultAppMode = "debug"
)

type WebServerConfig struct {
	Host string
	Port int64
	Mode string // https://github.com/gin-gonic/gin/blob/dc9cff732e27ce4ac21b25772a83c462a28b8b80/mode.go#L18
}

type AppConfig struct {
	PgConfig PgConfig
	WsConfig WebServerConfig
}

var Config AppConfig

func getEnvStrDefault(envVariables map[string]string, key string, def string) string {
	if value, exist := envVariables[key]; exist {
		return value
	}
	return def
}

func getEnvIntDefault(envVariables map[string]string, key string, def int64) int64 {
	if value, exist := envVariables[key]; exist {
		if result, err := strconv.ParseInt(value, 10, 64); err == nil {
			return result
		} else {
			fmt.Printf("Environment variable '%s' is not an integer (value '%s')", key, value)
			os.Exit(1)
		}
	}
	return def
}

func getEnvVariables() map[string]string {
	envVariables := make(map[string]string)
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			envVariables[e[:i]] = e[i+1:]
		}
	}
	return envVariables
}

func checkRequiredEnvVar(envVariables map[string]string) {
	var requiredEnvVar = [...]string{
		envPgHost,
		envPgUser,
		envPgPassword,
	}
	for _, envVar := range requiredEnvVar {
		if varValue, exist := envVariables[envVar]; !exist || len(varValue) == 0 {
			if exist {
				fmt.Printf("Environment variable '%s' value is empty", envVar)
				os.Exit(1)
			} else {
				fmt.Printf("Environment variable '%s' is required", envVar)
				os.Exit(1)
			}
		}
	}
}

func InitConfig() {
	envVariables := getEnvVariables()
	checkRequiredEnvVar(envVariables)

	config := AppConfig{
		PgConfig: PgConfig{
			PgHost:     envVariables[envPgHost],
			PgPort:     getEnvIntDefault(envVariables, envPgPort, defaultPgPort),
			PgDbName:   getEnvStrDefault(envVariables, envPgDB, defaultPgDb),
			PgUser:     envVariables[envPgUser],
			PgPassword: envVariables[envPgPassword],
			PgSslMode:  getEnvStrDefault(envVariables, envPgSslMode, defaultPgSslMode),
		},
		WsConfig: WebServerConfig{
			Host: getEnvStrDefault(envVariables, envAppHost, defaultAppHost),
			Port: getEnvIntDefault(envVariables, envAppPort, defaultAppPort),
			Mode: getEnvStrDefault(envVariables, envAppMode, defaultAppMode),
		},
	}
	Config = config
}
