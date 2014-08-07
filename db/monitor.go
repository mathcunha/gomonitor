package db

import "gopkg.in/mgo.v2/bson"
import "fmt"

type MonitorDB struct{
	Id		bson.ObjectId	`json:"id"		bson:"_id"`
	Query		string		`json:"query"		bson:"query"`
	Threshold	int		`json:"threshold"	bson:"threshold"`
	Interval	string		`json:"interval"	bson:"interval"`
	field		string		`json:"field"		bson:"field"`
}

func FindOneMonitor(id string) (error, MonitorDB){
	var monitor MonitorDB
	fmt.Println(id)
	err := FindOne("monitor",bson.M{"_id": id},&monitor)
	return err, monitor
}
