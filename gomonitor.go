package main

import (
	"github.com/mathcunha/gomonitor/handler"
	"github.com/mathcunha/gomonitor/scheduler"
	"github.com/mathcunha/gomonitor/prop"
	"net/http"
)

func main() {
	scheduler.LoadMonitors()
	http.HandleFunc("/", handler.DoRequest)
	http.ListenAndServe(prop.Property("gomonitor"), nil)
}
