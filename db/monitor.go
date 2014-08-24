package db

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

var collection = "monitor"

type Monitor struct {
	Id        bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Query     string        `bson:"query,omitempty" json:"query"`
	Threshold int           `bson:"threshold"       json:"threshold"`
	Interval  string        `bson:"interval"        json:"interval"`
	Field     string        `bson:"field"           json:"field"`
	Actions   []string      `bson:"actions"         json:"actions"`
}

func FindOneMonitor(id string) (error, Monitor) {
	var monitor Monitor
	err := FindOne(collection, bson.M{"_id": bson.ObjectIdHex(id)}, &monitor)
	return err, monitor
}

func FindAllMonitor() (error, []Monitor) {
	var monitors []Monitor
	err := FindAll(collection, &monitors)
	return err, monitors
}

func InsertMonitor(decoder *json.Decoder) (error, Monitor) {
	var monitor Monitor
	err := decoder.Decode(&monitor)

	if err != nil {
		return err, monitor
	}

	return Insert(collection, &monitor), monitor
}

func RemoveMonitor(id string) error {
	return Remove(collection, bson.M{"_id": bson.ObjectIdHex(id)})
}
