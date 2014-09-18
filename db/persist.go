package db

import (
	"github.com/mathcunha/gomonitor/prop"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Persistent interface {
	collection() string
}

const database = "gomonitor"

func getSession() (*mgo.Session, error) {
	session, err := mgo.Dial("mongodb://" + prop.Property("mongodb"))
	return session, err
}

func closeSession(s *mgo.Session) {
	s.Close()
}

func GetId(id string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id)}
}

func FindOne(document Persistent, id bson.M) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	return s.DB(database).C(document.collection()).Find(id).One(document)
}

func FindQuery(collection string, result interface{}, query interface{}) error {
	s, err := getSession()
	if err != nil {
		return err
	}

	defer closeSession(s)

	return s.DB(database).C(collection).Find(query).All(result)
}

func FindAll(collection string, result interface{}) error {
	return FindQuery(collection, result, nil)
}

func Insert(document Persistent) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	return s.DB(database).C(document.collection()).Insert(document)
}

func Remove(document Persistent, id bson.M) error {
	s, err := getSession()
	if err != nil {
		return err
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	return s.DB(database).C(document.collection()).Remove(id)
}
