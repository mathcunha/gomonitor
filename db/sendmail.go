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
	From    string        `bson:"from"              json:"from"`
}

func (s Sendmail) collection() string {
	return "sendmail"
}

func (s Sendmail) FindAll() (error, []Sendmail) {
	var sendmails []Sendmail
	err := FindAll("sendmail", &sendmails)
	return err, sendmails
}

func (s Sendmail) FindByMonitor(m Monitor) (error, Sendmail) {
	sendmail := new(Sendmail)
	return FindOne(sendmail, bson.M{"monitor._id": m.Id}), *sendmail
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
