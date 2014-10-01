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
