package main

import (
	"container/list"
	"errors"
	"log"
	"time"

	"github.com/xu354cjo1008/eatingFinder/meteorology"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct {
	databaseUrl string
	sessionMaxN int
	sessions    *list.List
}

type RestaurantInfo struct {
	name     string
	open_now string
	place_id string
	rating   float64
	vicinity string
}

type ChoiceElement struct {
	Lat        float64
	Lng        float64
	Time       time.Time
	Restaurant RestaurantInfo
	Weather    meteorology.Weather
}

func (storage *Storage) insertChoice(db *mgo.Database, element ChoiceElement) error {

	collection := db.C("restaurant_choice")
	err := collection.Insert(element)
	if err != nil {
		return err
	}
	return nil
}

func (storage *Storage) findChoiceListByLocation(db *mgo.Database, lat float64, lng float64, size int) []ChoiceElement {

	collection := db.C("restaurant_choice")

	countNum, err := collection.Count()
	if err != nil {
		log.Println(err)
	}

	log.Println("Things objects count: ", countNum)

	result := []ChoiceElement{}
	//	collection.Find(bson.M{"lat": lat}).All(&result)

	collection.Find(bson.M{
		"$and": []bson.M{
			bson.M{
				"lat": bson.M{
					"$gt": lat - 10,
					"$lt": lat + 10,
				},
			},
			bson.M{
				"lng": bson.M{
					"$gt": lng - 10,
					"$lt": lng + 10,
				},
			},
		},
	}).All(&result)

	return result
}

func (storage *Storage) getDb(name string, user string, password string) (*mgo.Database, error) {

	if storage.sessions.Len() > storage.sessionMaxN {
		return nil, errors.New("there is no free session can be used")
	}

	mgoSession, err := mgo.Dial(storage.databaseUrl)
	if err != nil {
		log.Println("mgoSession failed")
		return nil, err
	}

	db := mgoSession.DB(name)

	if db == nil {
		log.Println("empty db")
	}

	err = mgoSession.Login(&mgo.Credential{Username: user, Password: password, Source: name})
	if err != nil {
		log.Fatalln("can not login ", err)
		return nil, err
	}

	storage.sessions.PushBack(mgoSession)

	return db, nil
}

func (storage *Storage) close(session *mgo.Session) {

	for e := storage.sessions.Front(); e != nil; e = e.Next() {
		if e.Value.(*mgo.Session) == session {
			storage.sessions.Remove(e)
			e.Value.(*mgo.Session).Close()
		}
	}
}

func NewStorage(url string) *Storage {

	storage := Storage{
		databaseUrl: url,
		sessions:    list.New(),
	}
	return &storage
}
