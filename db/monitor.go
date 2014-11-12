package db

import (
	"github.com/mathcunha/gomonitor/es"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

type Monitor struct {
	Id        bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Query     string        `bson:"query,omitempty" json:"query"`
	Intervalo string        `bson:"interval"        json:"interval"`
	Actions   []string      `bson:"actions"         json:"actions"`
}

func (m Monitor) collection() string {
	return "monitor"
}

func (m Monitor) FindAll() ([]Monitor, error) {
	var monitors []Monitor
	err := FindAll(m.collection(), &monitors)
	return monitors, err
}

func (m Monitor) Interval() string {
	return m.Intervalo
}

func (m Monitor) Run() {
	total, hits, err := es.Search(strings.Replace(m.Query, "\\", "", -1))
	log.Printf("monitor = (%v) total = (%v) - hits = (%v)", m.Id, total, hits)

	if err == nil && total > 0 {
		var alert Alert
		alert.Monitor = m
		alert.Hits = hits
		alert.Total = total
		alert.When = time.Now()
		alert.PostAlert()
	}
}
