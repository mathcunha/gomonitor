package db

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

var collection = "monitor"

type Monitor struct {
	Id       bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Query    string        `bson:"query,omitempty" json:"query"`
	Interval string        `bson:"interval"        json:"interval"`
	Actions  []string      `bson:"actions"         json:"actions"`
}

func (m Monitor) FindOne(id string) (error, Monitor) {
	var monitor Monitor
	err := FindOne(collection, bson.M{"_id": bson.ObjectIdHex(id)}, &monitor)
	return err, monitor
}

func (m Monitor) FindAll() (error, []Monitor) {
	var monitors []Monitor
	err := FindAll(collection, &monitors)
	return err, monitors
}

func (m Monitor) Insert(decoder *json.Decoder) (error, Monitor) {
	var monitor Monitor
	err := decoder.Decode(&monitor)

	if err != nil {
		return err, monitor
	}

	return Insert(collection, &monitor), monitor
}

func (m Monitor) Remove(id string) error {
	return Remove(collection, bson.M{"_id": bson.ObjectIdHex(id)})
}
