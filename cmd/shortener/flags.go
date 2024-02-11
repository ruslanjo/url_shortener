package main

import (
	"flag"
	"github.com/ruslanjo/url_shortener/internal/config"
)


func parseFlags(){
	flag.StringVar(&config.ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.Protocol, "p", "http", "Protocol to run server")
	flag.StringVar(&config.BaseServerReturnAddr, "b", "http://localhost:8080", "Base addres of URL shortener")
	flag.Parse()
}


