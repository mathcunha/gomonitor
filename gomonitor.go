package main

import (
	"github.com/mathcunha/gomonitor/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/monitor/", handler.DoRequest)
	http.ListenAndServe(":8080", nil)
}
