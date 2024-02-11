package config

import (
	"flag"
	"os"
	"fmt"
)

func parseFlags() {
	flag.StringVar(&ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Protocol, "p", "http", "Protocol to run server")
	flag.StringVar(&BaseServerReturnAddr, "b", "http://localhost:8080", "Base addres of URL shortener")
	flag.Parse()
}

func setEnvsToConfig(variable *string, envName string) {
	if envValue, ok := os.LookupEnv(envName); ok {
		*variable = envValue
	}
}

func ConfigureApp() {
	parseFlags()

	envs := map[string]*string{
		"RUN_ADDR": &ServerAddr,
		"BASE_URL": &BaseServerReturnAddr,
	}

	for k, v := range envs {
		setEnvsToConfig(v, k)
	}
	fmt.Println(ServerAddr, BaseServerReturnAddr)
}
