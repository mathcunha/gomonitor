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
	monitor := new(db.Monitor)
	return db.FindOne(monitor, db.GetId(id)), monitor
}
func (m Monitor) insert(decoder *json.Decoder) (error, interface{}) {
	var monitor db.Monitor
	err := decoder.Decode(&monitor)

	if err != nil {
		return err, monitor
	}

	return db.Insert(monitor), monitor
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
	a := new(db.Alert)
	return db.FindOne(a, db.GetId(id)), a
}

func (alert Alert) insert(decoder *json.Decoder) (error, interface{}) {
	var a db.Alert
	err := decoder.Decode(&a)
	if err != nil {
		return err, a
	}

	err, alert_db := db.Insert(a), a

	for _, value := range alert_db.Monitor.Actions {
		log.Printf("posting alert to %v", value)
		var postData []byte
		w := bytes.NewBuffer(postData)
		json.NewEncoder(w).Encode(alert_db)
		http.Post("http://"+prop.Property("gomonitor")+"/"+value+"/action", "application/json", w)
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
	s := new(db.Sendmail)
	return db.FindOne(s, db.GetId(id)), s
}

func (sendmail Sendmail) insert(decoder *json.Decoder) (error, interface{}) {
	var s db.Sendmail
	log.Printf("sendmail.action = [%v]", sendmail.action)

	if "action" == sendmail.action {
		var alert db.Alert
		err := decoder.Decode(&alert)

		if err != nil {
			return err, alert
		}

		log.Printf("monitor %v", alert.Monitor)

		err, s = s.FindByMonitor(alert.Monitor)

		if err != nil {
			log.Printf("You must insert a sendmail to monitor %v ", alert.Monitor.Id)
			return err, alert
		}

		action.SimpleSendMail(s.From, s.To, fmt.Sprintf("%v", alert.Monitor.Id), fmt.Sprintf("%v", alert))

		return err, s
	} else {
		err := decoder.Decode(&s)
		if err != nil {
			return err, s
		}
		return db.Insert(s), s
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
