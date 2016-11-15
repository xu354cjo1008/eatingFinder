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
	Name     string
	Open_now string
	Place_id string
	Rating   float64
	Vicinity string
	Rank     int
}

type DiscoverInfo struct {
	Lat    float64
	Lng    float64
	Radius float64
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

func (storage *Storage) insertDiscoverInfo(db *mgo.Database, element DiscoverInfo) error {

	collection := db.C("restaurant_discover")
	err := collection.Insert(element)
	if err != nil {
		return err
	}
	return nil
}

func (storage *Storage) findDiscoverInfo(db *mgo.Database, lat float64, lng float64, maxRadius float64) ([]DiscoverInfo, error) {
	collection := db.C("restaurant_discover")

	countNum, err := collection.Count()
	if err != nil {
		return nil, err
	}
	var result []DiscoverInfo

	if countNum > 0 {
		collection.Find(bson.M{
			"$and": []bson.M{
				bson.M{
					"lat": bson.M{
						"$gt": lat - 1,
						"$lt": lat + 1,
					},
				},
				bson.M{
					"lng": bson.M{
						"$gt": lng - 1,
						"$lt": lng + 1,
					},
				},
			},
		}).All(&result)
	}

	return result, nil
}

func (storage *Storage) findChoiceListByLocation(db *mgo.Database, lat float64, lng float64, size float64) []ChoiceElement {

	collection := db.C("restaurant_choice")

	countNum, err := collection.Count()
	if err != nil {
		log.Println(err)
	}

	log.Println("Things objects count: ", countNum)

	result := []ChoiceElement{}

	collection.Find(bson.M{
		"$and": []bson.M{
			bson.M{
				"lat": bson.M{
					"$gt": lat - 1,
					"$lt": lat + 1,
				},
			},
			bson.M{
				"lng": bson.M{
					"$gt": lng - 1,
					"$lt": lng + 1,
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
