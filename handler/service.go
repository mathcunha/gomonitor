package handler

import (
	"fmt"
	"net/http"
	"strings"
	"../db"
)

type ServiceHandler interface{
	get(w http.ResponseWriter)
	getOne(w http.ResponseWriter, id string)
}

type Monitor struct{}

func getId(r *http.Request) string{
	data := strings.Split(r.URL.Path,"/")

	//for i, v := range data{
	//	fmt.Printf("indice = valor | %d = %d\n", i, v)
	//}

	if len(data) > 2{
		return data[2]
	}

	return ""
}

func getHandler(path string) ServiceHandler{
	a_path := strings.Split(path,"/")
	if "monitor" == a_path[1]{
		return Monitor{}
	}
	return nil
}

func DoRequest (w http.ResponseWriter, r *http.Request) {

	var monitorHandler ServiceHandler
	monitorHandler = getHandler(r.URL.Path)

	switch{
	case "GET" == r.Method:
		id := getId(r)
		fmt.Printf("temos um GET - %v \n", id)
		if(id != ""){
			monitorHandler.getOne(w, id)
		}else{
			monitorHandler.get(w)
		}
	case "DELETE" == r.Method:
		fmt.Printf("temos um DELETE - %v \n", getId(r))
	case "POST" == r.Method:
		fmt.Println("temos um POST")
	}
}

func (monitor Monitor) get (w http.ResponseWriter){
	fmt.Fprintf(w, "<h1>Editing %s</h1>", "GET")
}
func (monitor Monitor) getOne (w http.ResponseWriter, id string){
	err, result := db.FindOneMonitor(id)
	if err == nil{
        	fmt.Fprintf(w, "<h1>Editing %s</h1>", result)
	}else{
		fmt.Fprintf(w, "<h1>Editing %s</h1>", err)
	}
}
