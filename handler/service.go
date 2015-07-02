package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mathcunha/gomonitor/action"
	"github.com/mathcunha/gomonitor/db"
	"github.com/mathcunha/gomonitor/prop"
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
	getAll() (interface{}, error)
	getOne(id string) (interface{}, error)
	insert(decode *json.Decoder) (interface{}, error)
	removeOne(id string) error
}

func getId(r *http.Request) string {
	data := strings.Split(r.URL.Path, "/")
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
		if len(a_path) >= 3 {
			return Sendmail{a_path[2]}
		} else {
			return Sendmail{a_path[1]}
		}
	} else if "alert" == a_path[1] {
		return Alert{}
	}
	return nil
}

func DoRequest(w http.ResponseWriter, r *http.Request) {

	var monitorHandler ServiceHandler

	if monitorHandler = getHandler(r.URL.Path); monitorHandler == nil {
		http.Error(w, "no handler to path "+r.URL.Path, http.StatusNotFound)
		log.Printf("no handler to path [%v]", r.URL.Path)
		return
	}

	var err error
	var result interface{}

	switch {
	case "GET" == r.Method:
		id := getId(r)
		if id != "" {
			result, err = monitorHandler.getOne(id)
		} else {
			result, err = monitorHandler.getAll()
		}
	case "DELETE" == r.Method:
		err = monitorHandler.removeOne(getId(r))
	case "POST" == r.Method:
		decoder := json.NewDecoder(r.Body)
		result, err = monitorHandler.insert(decoder)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		log.Printf("error handling %q: %v", r.RequestURI, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (m Monitor) getOne(id string) (interface{}, error) {
	monitor := new(db.Monitor)
	return monitor, db.FindOne(monitor, db.GetId(id))
}
func (m Monitor) insert(decoder *json.Decoder) (interface{}, error) {
	var monitor db.Monitor
	err := decoder.Decode(&monitor)

	if err != nil {
		return monitor, err
	}

	return monitor, db.Insert(monitor)
}

func (m Monitor) getAll() (interface{}, error) {
	var monitor db.Monitor
	return monitor.FindAll()
}

func (m Monitor) removeOne(id string) error {
	var monitor db.Monitor
	return db.Remove(monitor, db.GetId(id))
}

func (alert Alert) getOne(id string) (interface{}, error) {
	a := new(db.Alert)
	return a, db.FindOne(a, db.GetId(id))
}

func (alert Alert) insert(decoder *json.Decoder) (interface{}, error) {
	var a db.Alert
	err := decoder.Decode(&a)
	if err != nil {
		return a, err
	}

	err, alert_db := db.Insert(a), a

	for _, value := range alert_db.Monitor.Actions {
		log.Printf("posting alert to %v", value)
		var postData []byte
		w := bytes.NewBuffer(postData)
		json.NewEncoder(w).Encode(alert_db)
		resp, err := http.Post("http://"+prop.Property("gomonitor")+"/"+value+"/action", "application/json", w)
		if err != nil {
			return a, err
		}
		defer resp.Body.Close()
	}

	return a, nil
}

func (alert Alert) getAll() (interface{}, error) {
	var a db.Alert
	return a.FindAll()
}

func (alert Alert) removeOne(id string) error {
	var a db.Alert
	return db.Remove(a, db.GetId(id))
}

func (sendmail Sendmail) getOne(id string) (interface{}, error) {
	s := new(db.Sendmail)
	return s, db.FindOne(s, db.GetId(id))
}

func (sendmail Sendmail) insert(decoder *json.Decoder) (interface{}, error) {
	var s db.Sendmail
	log.Printf("sendmail.action = [%v]", sendmail.action)

	if "action" == sendmail.action {
		var alert db.Alert
		err := decoder.Decode(&alert)

		if err != nil {
			return alert, err
		}

		log.Printf("monitor %v", alert.Monitor)

		err, s = s.FindByMonitor(alert.Monitor)

		if err != nil {
			log.Printf("You must insert a sendmail to monitor %v ", alert.Monitor.Id)
			return alert, err
		}

		action.SimpleSendMail(s.From, s.To, fmt.Sprintf("%v", alert.Monitor.Id), fmt.Sprintf("%v", alert))

		return alert, err
	} else {
		err := decoder.Decode(&s)
		if err != nil {
			return s, err
		}
		return s, db.Insert(s)
	}
}

func (sendmail Sendmail) getAll() (interface{}, error) {
	var s db.Sendmail
	return s.FindAll()
}

func (sendmail Sendmail) removeOne(id string) error {
	var s db.Sendmail
	return db.Remove(s, db.GetId(id))
}
