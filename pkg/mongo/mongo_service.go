package mongo

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
)

func newMongoSession() *mgo.Session {
	//TODO - load from VCAP_SERVICES, etc
	mongoDialInfo := &mgo.DialInfo{
		Addrs:    []string{"127.0.0.1"},
		Database: "cohesion-education",
		Username: "admin",
		Password: "password",
		Source:   "admin",
		Timeout:  60 * time.Second,
	}
	session, err := mgo.DialWithInfo(mongoDialInfo)
	if err != nil {
		log.Fatalf("Failed to connect to mongodb: %v", err)
	}

	return session
}
