package main

import (
	"flag"
	"github.com/mathcunha/amon/scheduler"
	"github.com/mathcunha/gomonitor/db"
	"github.com/mathcunha/gomonitor/handler"
	"github.com/mathcunha/gomonitor/prop"
	"log"
	"net/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "prop/config.json", "gomonitor configuration file")
}

func main() {
	flag.Parse()
	prop.LoadConfig(configFile)
	LoadMonitors()
	http.HandleFunc("/", handler.DoRequest)
	http.ListenAndServe(prop.Property("gomonitor"), nil)
}

func LoadMonitors() error {
	var monitor db.Monitor
	monitors, err := monitor.FindAll()

	if err != nil {
		log.Printf("error loading the montitors %v", err)
		return err
	}
	scheduler.Schedule(tasks(monitors))

	return nil
}

func tasks(s []db.Monitor) []scheduler.Task {
	vals := make([]scheduler.Task, len(s))
	for i, v := range s {
		vals[i] = v
	}
	return vals
}
