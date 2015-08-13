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
	Hits      int           `bson:"hits"         json:"hits"`
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
	hits, err := es.Search(strings.Replace(m.Query, "\\", "", -1))
	log.Printf("monitor = (%v) - hits = (%v)", m.Id, hits)

	if err == nil {
		for _, h := range hits {
			if h.Total > m.Hits {
				var alert Alert
				alert.Monitor = m
				alert.Hits = h.Prop
				alert.Total = h.Total
				alert.When = time.Now()
				alert.PostAlert()
			}
		}
	}
}
