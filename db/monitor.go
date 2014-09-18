package db

import (
	"gopkg.in/mgo.v2/bson"
)

type Monitor struct {
	Id       bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Query    string        `bson:"query,omitempty" json:"query"`
	Interval string        `bson:"interval"        json:"interval"`
	Actions  []string      `bson:"actions"         json:"actions"`
}

func (m Monitor) collection() string {
	return "monitor"
}

func (m Monitor) FindAll() ([]Monitor, error) {
	var monitors []Monitor
	err := FindAll(m.collection(), &monitors)
	return monitors, err
}
