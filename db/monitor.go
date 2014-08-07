package db

import "gopkg.in/mgo.v2/bson"

type Monitor struct{
	Id		bson.ObjectId	"_id,omitempty"
	Query		string		"query"
	Threshold	int
	Interval	string
	field		string
}

func FindOneMonitor(id string) (error, Monitor){
	var monitor Monitor
	err := FindOne("monitor",bson.M{"_id": bson.ObjectIdHex(id)}, &monitor)
	return err, monitor
}
