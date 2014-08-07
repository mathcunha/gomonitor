package main

import (
	"net/http"
	"./handler"
)

func main(){
	http.HandleFunc("/monitor/", handler.DoRequest)
	http.ListenAndServe(":8080", nil)
}
