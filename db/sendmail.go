package db

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Sendmail struct {
	Id      bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Monitor Monitor       `bson:"monitor"      json:"monitor"`
	To      []string      `bson:"to"              json:"to"`
}

func (s Sendmail) FindOne(id string) (error, Sendmail) {
	var sendmail Sendmail
	err := FindOne("sendmail", bson.M{"_id": bson.ObjectIdHex(id)}, &sendmail)
	return err, sendmail
}

func (s Sendmail) FindAll() (error, []Sendmail) {
	var sendmails []Sendmail
	err := FindAll("sendmail", &sendmails)
	return err, sendmails
}

func (s Sendmail) Insert(decoder *json.Decoder) (error, Sendmail) {
	var sendmail Sendmail
	err := decoder.Decode(&sendmail)

	if err != nil {
		return err, sendmail
	}

	return Insert("sendmail", &sendmail), sendmail
}

func (s Sendmail) Action(decoder *json.Decoder) (error, Alert) {
	var alert Alert
	err := decoder.Decode(&alert)

	if err != nil {
		return err, alert
	}

	//TODO sendmail stuff
	log.Printf("sending email - [%v]", alert.Monitor.Id)

	return err, alert
}

func (s Sendmail) Remove(id string) error {
	return Remove("sendmail", bson.M{"_id": bson.ObjectIdHex(id)})
}
