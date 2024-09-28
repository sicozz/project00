package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sicozz/project00/utils"
)

const (
	DEFAULT_HOST = "[::]"
	DEFAULT_PORT = "50050"
)

type Config struct {
	Host      string
	Port      string
	HostsFile string
	LogFile   string
}

func BuildConfig() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load .env file")
	}
	host := getEnvElseDefault("HOST", DEFAULT_HOST)
	port := getEnvElseDefault("PORT", DEFAULT_PORT)
	logFile := getEnvElseDefault("LOG_FILE", utils.DEFAULT_LOG_FILE)
	hostsFile := getEnvElseDefault("HOSTS_FILE", utils.DEFAULT_HOSTS_FILE)

	return Config{
		Host:      host,
		Port:      port,
		LogFile:   logFile,
		HostsFile: hostsFile,
	}
}

func (c *Config) GetBindAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnvElseDefault(varName, defaultValue string) string {
	envVar := os.Getenv(varName)
	if envVar == "" {
		return defaultValue
	}
	return envVar
}
