package db

import (
	"labix.org/v2/mgo"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func getSession() (*mgo.Session, error){
	session, err := mgo.Dial("mongodb://127.0.0.1")
	return session, err
}

func closeSession(s *mgo.Session){
	s.Close()
}

func FindOne(collection string, id bson.M, result *MonitorDB) error{
	s, err := getSession()
	fmt.Println(err)
	fmt.Println(id)
	err = s.DB("gomonitor").C(collection).FindId(id).One(result)
	fmt.Println(err)
	closeSession(s)
	return err
}
