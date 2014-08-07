package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getSession() (*mgo.Session, error){
	session, err := mgo.Dial("mongodb://127.0.0.1")
	return session, err
}

func closeSession(s *mgo.Session){
	s.Close()
}

func FindOne(collection string, id bson.M, result interface{}) error{
	s, err := getSession()
	if err != nil {
		panic(err)
	}
	defer closeSession(s)

	err = s.DB("gomonitor").C(collection).Find(id).One(result)

	return err
}
