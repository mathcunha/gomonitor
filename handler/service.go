package handler

import (
	"fmt"
	"net/http"
	"strings"
	"../db"
	"encoding/json"
)
type Monitor struct{}
type ServiceHandler interface{
	getAll() (error, interface{})
	getOne(id string) (error, interface{})
	insert(decode *json.Decoder) (error, interface{})
	removeOne (id string) (error)
}

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

	var err error
	var result interface{}

	switch{
	case "GET" == r.Method:
		id := getId(r)
		if(id != ""){
			err, result = monitorHandler.getOne(id)
		}else{
			err, result = monitorHandler.getAll()
		}
	case "DELETE" == r.Method:
		err = monitorHandler.removeOne(getId(r))
	case "POST" == r.Method:
		decoder := json.NewDecoder(r.Body)
		err, result = monitorHandler.insert(decoder)
	}

	if err == nil{
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}else{
		fmt.Fprintf(w, "error %s", err)
	}
}

func (monitor Monitor) getOne (id string) (error, interface{}){
	return db.FindOneMonitor(id)
}
func (monitor Monitor) insert (decoder *json.Decoder)(error, interface{}){
	return db.InsertMonitor(decoder)
}

func (monitor Monitor) getAll()(error, interface{}){
	return db.FindAllMonitor()
}

func (monitor Monitor) removeOne (id string) (error){
	return db.RemoveMonitor(id)
}
