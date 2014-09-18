package db

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Alert struct {
	Id      bson.ObjectId `bson:"_id,omitempty"   json:"id,omitempty"`
	Monitor Monitor       `bson:"monitor"         json:"monitor"`
	Total   int           `bson:"total"        json:"total"`
	When    time.Time     `bson:"when"            json:"when"`
}

func (a Alert) collection() string {
	return "alert"
}

func (a Alert) FindAll() (error, []Alert) {
	var alerts []Alert
	err := FindAll("alert", &alerts)
	return err, alerts
}

func (a Alert) Remove(id string) error {
	return Remove("alert", bson.M{"_id": bson.ObjectIdHex(id)})
}
