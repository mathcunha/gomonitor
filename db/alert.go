package db

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type Alert struct {
	Id       bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Monitor  Monitor       `bson:"monitor"         json:"monitor"`
	Measured int           `bson:"threshold"       json:"threshold"`
}

func FindOneAlert(id string) (error, Alert) {
	var alert Alert
	err := FindOne("alert", bson.M{"_id": bson.ObjectIdHex(id)}, &alert)
	return err, alert
}

func FindAllAlert() (error, []Alert) {
	var alerts []Alert
	err := FindAll("alert", &alerts)
	return err, alerts
}

func InsertAlert(decoder *json.Decoder) (error, Alert) {
	var alert Alert
	err := decoder.Decode(&alert)

	if err != nil {
		return err, alert
	}

	return Insert("alert", &alert), alert
}

func RemoveAlert(id string) error {
	return Remove("alert", bson.M{"_id": bson.ObjectIdHex(id)})
}
