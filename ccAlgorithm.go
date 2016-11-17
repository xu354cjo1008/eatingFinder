package main

import (
	"io"
	"log"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"github.com/kr/pretty"
	"github.com/xu354cjo1008/eatingFinder/geography/place"
	"github.com/xu354cjo1008/eatingFinder/meteorology"
)

type ccAlgorithm struct {
	place    *nearPlace.GoogleBase
	meteo    *meteorology.Meteorology
	storage  *Storage
	logLevel int
	logger   *log.Logger
}

func (alg *ccAlgorithm) checkIsDiscovered(db *mgo.Database, lat float64, lng float64, size float64) bool {

	elements, err := alg.storage.findDiscoverInfo(db, lat, lng, size)
	if err != nil {
		if alg.logLevel == 1 {
			alg.logger.Println(err)
		}
		return false
	}
	for _, element := range elements {
		//now just simple check by rectangle for test
		if lat-size >= element.Lat-element.Radius && lat+size <= element.Lat+element.Radius {
			if lng-size >= element.Lng-element.Radius && lng+size <= element.Lng+element.Radius {
				if alg.logLevel == 1 {
					alg.logger.Println("discovered areas")
					alg.logger.Println(elements)
				}
				return true
			}
		}
	}
	return false
}

func (alg *ccAlgorithm) findRestaurant(lat float64, lng float64) {

}

func (alg *ccAlgorithm) findRestaurantList(mode int, userData algUserData, size int) {

	var err error
	var db *mgo.Database
	if alg.logLevel == 1 {
		alg.logger.Println("enter findRestaurantList -> lat: ", userData.lat, "lng: ", userData.lng)
	}
	// try to get db instance(maybe failed because there are no enougth session in pool)
	db, err = alg.storage.getDb(config.dbName, config.dbUsername, config.dbPassword)
	if err != nil {
		if alg.logLevel == 1 {
			alg.logger.Println(err)
		}
		return
	}
	isDiscovered := alg.checkIsDiscovered(db, userData.lat, userData.lng, float64(size))
	// search data from storage
	if isDiscovered {
		choice := alg.storage.findChoiceListByLocation(db, userData.lat, userData.lng, float64(size))
		if alg.logLevel == 1 {
			alg.logger.Println("The area has been discoverd")
			alg.logger.Println(choice)
		}
	} else {
		// and then query from remote api
		data, err := alg.place.GetNearRestaurants(userData.lat, userData.lng, uint(size), "en")
		if err != nil {
			if alg.logLevel == 1 {
				alg.logger.Println(err)
			}
			return
		}
		switch mode {
		case ALG_HIGHEST_RATE:
			if alg.logLevel == 1 {
				pretty.Println(data)
				for _, paper := range data {
					for rank, info := range paper {
						restaurantElement := RestaurantInfo{}
						element := ChoiceElement{}
						for field, value := range info {
							switch field {
							case "vicinity":
								restaurantElement.Vicinity = value.(string)
							case "Location":
								s := strings.Split(value.(string), ",")
								lat := strings.Split(s[0], ":")
								lng := strings.Split(s[1], ":")
								lat64, err := strconv.ParseFloat(strings.TrimSpace(lat[1]), 64)
								if err != nil {
									alg.logger.Println(err)
								}
								lng64, err := strconv.ParseFloat(strings.TrimSpace(lng[1]), 64)
								if err != nil {
									alg.logger.Println(err)
								}
								element.Lat = lat64
								element.Lng = lng64
							case "name":
								restaurantElement.Name = value.(string)
							case "place_id":
								restaurantElement.Place_id = value.(string)
							}
						}
						restaurantElement.Rank = rank
						element.Restaurant = restaurantElement

						err := alg.storage.insertChoice(db, element)
						if err != nil {
							alg.logger.Println(err)
						}
					}
				}
			}
		}

		err = alg.storage.insertDiscoverInfo(db, DiscoverInfo{Lat: userData.lat, Lng: userData.lng, Radius: float64(size)})
		if err != nil {
			if alg.logLevel == 1 {
				alg.logger.Println(err)
			}
			return
		}
	}

	alg.storage.close(db.Session)
}

func newCCAlgorithm(logFile io.Writer) *ccAlgorithm {

	var loggingLevel int
	if logFile == nil {
		loggingLevel = 0
	} else {
		loggingLevel = 1
	}

	near, err := nearPlace.InitPlaceNearbySearch("google_lib")
	if err != nil {
		log.Fatalln("Failed to create nearPlace instance")
		return nil
	}

	storage := NewStorage(config.dbUrl)

	alg := ccAlgorithm{
		place:    near,
		meteo:    meteorology.NewMeteorology(config.cwdApiKey, "en", nil),
		storage:  storage,
		logLevel: loggingLevel,
		logger: log.New(logFile, "ccAlgorithm: ",
			log.Ldate|log.Ltime|log.Lshortfile),
	}
	return &alg
}
