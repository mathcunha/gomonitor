package main

import (
	"github.com/mathcunha/gomonitor/handler"
	"github.com/mathcunha/gomonitor/scheduler"
	"net/http"
)

func main() {
	scheduler.LoadMonitors()
	http.HandleFunc("/", handler.DoRequest)
	http.ListenAndServe(":8080", nil)
}
