package db

import (
	"gopkg.in/mgo.v2/bson"
)

var collection = "monitor"

type Monitor struct {
	Id       bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Query    string        `bson:"query,omitempty" json:"query"`
	Interval string        `bson:"interval"        json:"interval"`
	Actions  []string      `bson:"actions"         json:"actions"`
}

func (m Monitor) collection() string {
	return "monitor"
}

func (m Monitor) FindAll() (error, []Monitor) {
	var monitors []Monitor
	err := FindAll(collection, &monitors)
	return err, monitors
}

func (m Monitor) Remove(id string) error {
	return Remove(collection, bson.M{"_id": bson.ObjectIdHex(id)})
}
