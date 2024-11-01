package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	Date    = ""
	Version = "dev"
	Build   = "dev"
)

const (
	EnvXeroAPISvrPort     = "XEROAPI_SVR_PORT"
	EnvXeroAPISvrLogLevel = "XEROAPI_SVR_LOG_LEVEL"
	EnvXeroAPISvrGinMode  = "XEROAPI_SVR_GIN_MODE"

	EnvFile = ".env"
)

func init() {
	fmt.Printf("Build Date: %s\nBuild Version: %s\nBuild: %s\n\n", Date, Version, Build)
	err := godotenv.Load(EnvFile)
	if err != nil {
		log.Fatalf("Error loading %s file: %v", EnvFile, err)
	}
	logLevel, err := log.ParseLevel(os.Getenv(EnvXeroAPISvrLogLevel))
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&log.JSONFormatter{})
}
