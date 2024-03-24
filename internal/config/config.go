package config

import (
	"fmt"
	"time"
)

const (
	DefaultServerHost string = "localhost"
	DefaultServerPort int    = 8080
)

var (
	Protocol             string = "http"
	ServerAddr           string = fmt.Sprintf("%s:%d", DefaultServerHost, DefaultServerPort)
	BaseServerReturnAddr string = fmt.Sprintf("%s://%s", Protocol, ServerAddr)
	LocalStoragePath     string
	DSN                  string
)

const (
	GZIP string = "gzip"
)

const (
	SelectedCompressionType string = GZIP
)

const (
	URLBatchSize int = 1_000
)

const (
	JWTSecret     string        = "s3cr3t" // it should be read from envs
	TokenLifeTime time.Duration = 3 * time.Hour
)
