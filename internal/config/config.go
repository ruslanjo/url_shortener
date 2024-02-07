package config

import "fmt"

var host string = "localhost"
var port int16 = 8080
var Addr string = fmt.Sprintf("%s:%d", host, port)
