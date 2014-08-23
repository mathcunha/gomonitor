package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const database = "gomonitor"

func getSession() (*mgo.Session, error) {
	session, err := mgo.Dial("mongodb://127.0.0.1")
	return session, err
}

func closeSession(s *mgo.Session) {
	s.Close()
}

func FindOne(collection string, id bson.M, result interface{}) error {
	s, err := getSession()
	if err != nil {
		panic(err)
	}
	defer closeSession(s)

	err = s.DB(database).C(collection).Find(id).One(result)

	return err
}

func FindQuery(collection string, result interface{}, query interface{}) error {
	s, err := getSession()
	if err != nil {
		panic(err)
	}

	defer closeSession(s)

	err = s.DB(database).C(collection).Find(query).All(result)

	return err
}

func FindAll(collection string, result interface{}) error {
	return FindQuery(collection, result, nil)
}

func Insert(collection string, document interface{}) error {
	s, err := getSession()
	if err != nil {
		panic(err)
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	err = s.DB(database).C(collection).Insert(document)

	return err
}

func Remove(collection string, id bson.M) error {
	s, err := getSession()
	if err != nil {
		panic(err)
	}
	defer closeSession(s)

	s.SetSafe(&mgo.Safe{FSync: true})

	err = s.DB(database).C(collection).Remove(id)

	return err
}
