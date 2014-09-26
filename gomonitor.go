package main

import (
	"flag"
	"github.com/mathcunha/gomonitor/handler"
	"github.com/mathcunha/gomonitor/prop"
	"github.com/mathcunha/gomonitor/scheduler"
	"net/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "prop/config.json", "gomonitor configuration file")
}

func main() {
	flag.Parse()
	prop.LoadConfig(configFile)
	scheduler.LoadMonitors()
	http.HandleFunc("/", handler.DoRequest)
	http.ListenAndServe(prop.Property("gomonitor"), nil)
}
