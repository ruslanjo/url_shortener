package config

import "fmt"

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
