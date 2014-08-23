package db

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type Sendmail struct {
	Id      bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Monitor Monitor       `bson:"monitor"      json:"monitor"`
	To      []string      `bson:"to"              json:"to"`
}

func FindOneSendmail(id string) (error, Sendmail) {
	var sendmail Sendmail
	err := FindOne("sendmail", bson.M{"_id": bson.ObjectIdHex(id)}, &sendmail)
	return err, sendmail
}

func FindAllSendmail() (error, []Sendmail) {
	var sendmails []Sendmail
	err := FindAll("sendmail", &sendmails)
	return err, sendmails
}

func InsertSendmail(decoder *json.Decoder) (error, Sendmail) {
	var sendmail Sendmail
	err := decoder.Decode(&sendmail)

	if err != nil {
		return err, sendmail
	}

	return Insert("sendmail", &sendmail), sendmail
}

func RemoveSendmail(id string) error {
	return Remove("sendmail", bson.M{"_id": bson.ObjectIdHex(id)})
}
