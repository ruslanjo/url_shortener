package config

import "fmt"

var Protocol string = "http"
var host string = "localhost"
var port int16 = 8080
var Addr string = fmt.Sprintf("%s:%d", host, port)
var FullAddr string = fmt.Sprintf("%s://%s", Protocol, Addr)


