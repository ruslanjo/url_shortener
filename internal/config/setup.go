package config

import (
	"flag"
	"os"
)

func parseFlags() {
	flag.StringVar(&ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Protocol, "p", "http", "Protocol to run server")
	flag.StringVar(&BaseServerReturnAddr, "b", "http://localhost:8080", "Base addres of URL shortener")
	flag.StringVar(&LocalStoragePath, "f", "/tmp/short-url-db.json", "Local DB addr")
	flag.Parse()
}

func setEnvsToConfig(envName string, variable *string) {
	if envValue, ok := os.LookupEnv(envName); ok {
		*variable = envValue
	}
}

func ConfigureApp() {
	parseFlags()

	envs := map[string]*string{
		"RUN_ADDR":          &ServerAddr,
		"BASE_URL":          &BaseServerReturnAddr,
		"FILE_STORAGE_PATH": &LocalStoragePath,
	}

	for k, v := range envs {
		setEnvsToConfig(k, v)
	}
}
