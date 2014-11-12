package db

import (
	"bytes"
	"encoding/json"
	"github.com/mathcunha/gomonitor/prop"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type Alert struct {
	Id      bson.ObjectId `bson:"_id,omitempty"   json:"id,omitempty"`
	Monitor Monitor       `bson:"monitor"         json:"monitor"`
	Total   int           `bson:"total"        json:"total"`
	When    time.Time     `bson:"when"            json:"when"`
	Hits    string        `bson:"hits"           json:"hits"`
}

func (a Alert) collection() string {
	return "alert"
}

func (a Alert) FindAll() ([]Alert, error) {
	var alerts []Alert
	err := FindAll(a.collection(), &alerts)
	return alerts, err
}

func (alert Alert) PostAlert() {
	var postData []byte
	w := bytes.NewBuffer(postData)
	json.NewEncoder(w).Encode(alert)
	http.Post("http://"+prop.Property("gomonitor")+"/alert", "application/json", w)
}
