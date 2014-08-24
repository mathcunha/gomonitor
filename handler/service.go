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

func (monitor Monitor) getOne(id string) (error, interface{}) {
	return db.FindOneMonitor(id)
}
func (monitor Monitor) insert(decoder *json.Decoder) (error, interface{}) {
	return db.InsertMonitor(decoder)
}

func (monitor Monitor) getAll() (error, interface{}) {
	return db.FindAllMonitor()
}

func (monitor Monitor) removeOne(id string) error {
	return db.RemoveMonitor(id)
}

func (alert Alert) getOne(id string) (error, interface{}) {
	return db.FindOneAlert(id)
}
func (alert Alert) insert(decoder *json.Decoder) (error, interface{}) {

	err, alert_db := db.InsertAlert(decoder)

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
	return db.FindAllAlert()
}

func (alert Alert) removeOne(id string) error {
	return db.RemoveAlert(id)
}

func (sendmail Sendmail) getOne(id string) (error, interface{}) {
	return db.FindOneSendmail(id)
}
func (sendmail Sendmail) insert(decoder *json.Decoder) (error, interface{}) {
	log.Printf("sendmail.action = [%v]", sendmail.action)

	if "action" == sendmail.action {
		return db.ActionSendmail(decoder)
	} else {
		return db.InsertSendmail(decoder)
	}
}

func (sendmail Sendmail) getAll() (error, interface{}) {
	return db.FindAllSendmail()
}

func (sendmail Sendmail) removeOne(id string) error {
	return db.RemoveSendmail(id)
}
