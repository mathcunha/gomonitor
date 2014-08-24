package handler

import (
	"bytes"
	"encoding/json"
	"github.com/mathcunha/gomonitor/db"
	"log"
	"net/http"
	"strings"
)

type Monitor struct{}

type Sendmail Action

type Action struct {
	action string
}

type Alert struct{}

type ServiceHandler interface {
	getAll() (error, interface{})
	getOne(id string) (error, interface{})
	insert(decode *json.Decoder) (error, interface{})
	removeOne(id string) error
}

func getId(r *http.Request) string {
	data := strings.Split(r.URL.Path, "/")

	//for i, v := range data{
	//	fmt.Printf("indice = valor | %d = %d\n", i, v)
	//}

	if len(data) > 2 {
		return data[2]
	}

	return ""
}

func getHandler(path string) ServiceHandler {
	a_path := strings.Split(path, "/")

	if "monitor" == a_path[1] {
		return Monitor{}
	} else if "sendmail" == a_path[1] {
		return Sendmail{a_path[2]}
	} else if "alert" == a_path[1] {
		return Alert{}
	}
	return nil
}

func DoRequest(w http.ResponseWriter, r *http.Request) {

	var monitorHandler ServiceHandler

	if monitorHandler = getHandler(r.URL.Path); monitorHandler == nil {
		http.Error(w, "no handler", http.StatusNotFound)
		log.Printf("1 - error handling %q", r.RequestURI)
		return
	}

	var err error
	var result interface{}

	switch {
	case "GET" == r.Method:
		id := getId(r)
		if id != "" {
			err, result = monitorHandler.getOne(id)
		} else {
			err, result = monitorHandler.getAll()
		}
	case "DELETE" == r.Method:
		err = monitorHandler.removeOne(getId(r))
	case "POST" == r.Method:
		decoder := json.NewDecoder(r.Body)
		err, result = monitorHandler.insert(decoder)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("2 - error handling %q: %v", r.RequestURI, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (m Monitor) getOne(id string) (error, interface{}) {
	var monitor db.Monitor
	return monitor.FindOne(id)
}
func (m Monitor) insert(decoder *json.Decoder) (error, interface{}) {
	var monitor db.Monitor
	return monitor.Insert(decoder)
}

func (m Monitor) getAll() (error, interface{}) {
	var monitor db.Monitor
	return monitor.FindAll()
}

func (m Monitor) removeOne(id string) error {
	var monitor db.Monitor
	return monitor.Remove(id)
}

func (alert Alert) getOne(id string) (error, interface{}) {
	var a db.Alert
	return a.FindOne(id)
}
func (alert Alert) insert(decoder *json.Decoder) (error, interface{}) {
	var a db.Alert
	err, alert_db := a.Insert(decoder)

	for _, value := range alert_db.Monitor.Actions {
		log.Printf("posting alert to %v", value)
		var postData []byte
		w := bytes.NewBuffer(postData)
		json.NewEncoder(w).Encode(alert_db)
		http.Post("http://127.0.0.1:8080/"+value+"/action", "application/json", w)
	}

	return err, alert
}

func (alert Alert) getAll() (error, interface{}) {
	var a db.Alert
	return a.FindAll()
}

func (alert Alert) removeOne(id string) error {
	var a db.Alert
	return a.Remove(id)
}

func (sendmail Sendmail) getOne(id string) (error, interface{}) {
	var s db.Sendmail
	return s.FindOne(id)
}
func (sendmail Sendmail) insert(decoder *json.Decoder) (error, interface{}) {
	var s db.Sendmail
	log.Printf("sendmail.action = [%v]", sendmail.action)

	if "action" == sendmail.action {
		return s.Action(decoder)
	} else {
		return s.Insert(decoder)
	}
}

func (sendmail Sendmail) getAll() (error, interface{}) {
	var s db.Sendmail
	return s.FindAll()
}

func (sendmail Sendmail) removeOne(id string) error {
	var s db.Sendmail
	return s.Remove(id)
}
