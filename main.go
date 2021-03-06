package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"github.com/xu354cjo1008/eatingFinder/geography/geocoding"
	"github.com/xu354cjo1008/eatingFinder/meteorology"
)

var config struct {
	defaultPort  int
	apiHost      string
	apiPort      int
	googleApiKey string
	cwdApiKey    string
	dbUrl        string
	dbName       string
	dbUsername   string
	dbPassword   string
}

func configure() error {

	viper.SetConfigName("app")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("Config file not found...")
	} else {
		config.defaultPort = viper.GetInt("development.defaultPort")
		config.apiHost = viper.GetString("development.apiHost")
		config.apiPort = viper.GetInt("development.apiPort")
		config.googleApiKey = viper.GetString("development.googleApiKey")
		config.cwdApiKey = viper.GetString("development.cwdApiKey")
		config.dbUrl = viper.GetString("development.dbUrl")
		config.dbName = viper.GetString("development.dbName")
		config.dbUsername = viper.GetString("development.dbUsername")
		config.dbPassword = viper.GetString("development.dbPassword")
	}

	log.Printf("\nDevelopment Config found:\n default server port = %d\n"+
		" api host = %s\n"+
		" api port = %d\n"+
		" db url = %s\n"+
		" db name = %s\n"+
		" db user = %s\n"+
		" db password = %s\n",
		config.defaultPort,
		config.apiHost,
		config.apiPort,
		config.dbUrl,
		config.dbName,
		config.dbUsername,
		config.dbPassword)

	return nil
}

func meteoUtil(lat float64, lng float64, logFile string) error {

	var file io.Writer = nil
	var err error

	if logFile != "" {
		if strings.Compare(logFile, "fg") == 0 {
			file = os.Stdout
		} else {
			file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalln("Failed to open log file :", err)
			}
		}
	}

	geocode := geocoding.NewGeocode(config.googleApiKey, "en")

	city, err := geocode.GetCityByLatlng(lat, lng)

	if err != nil {
		fmt.Println("error: ", err)
		return err
	}

	pretty.Println(city)

	meteo := meteorology.NewMeteorology(config.cwdApiKey, "en", file)
	data, err := meteo.GetWeather(city)
	if err != nil {
		log.Println("error: ", err)
		return err
	}

	pretty.Println(data)

	return nil
}

func algUtil(lat float64, lng float64, logFile string) error {

	var file io.Writer = nil
	var err error

	if logFile != "" {
		if strings.Compare(logFile, "fg") == 0 {
			file = os.Stdout
		} else {
			file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalln("Failed to open log file :", err)
			}
		}
	}

	alg := NewAlgorithm(file)
	alg.findRestaurantList(ALG_HIGHEST_RATE, algUserData{lat, lng}, 200)

	return nil
}

/**
 * This is the main just for test
 * We need to write another unit test program to do this
 */
func main() {

	var err error

	mode := flag.String("mode", "meteo", "utility mode: <meteo|web|api|load|save>")
	latPtr := flag.Float64("lat", 25.057339, "latitude of user position")
	lngPtr := flag.Float64("lng", 121.56086, "longtitude of user position")
	logFilePtr := flag.String("log", "", "log path <path|fg>")
	port := flag.Int("port", 0, "port number")

	flag.Parse()

	err = configure()

	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	if (strings.Compare(*mode, "web") == 0 || strings.Compare(*mode, "api") == 0) && *port != 0 {
		config.defaultPort = *port
	}

	switch *mode {
	case "meteo":
		err = meteoUtil(*latPtr, *lngPtr, *logFilePtr)
		if err != nil {
			pretty.Println(err)
			os.Exit(-1)
		}
	case "alg":
		err = algUtil(*latPtr, *lngPtr, *logFilePtr)
		if err != nil {
			pretty.Println(err)
			os.Exit(-1)
		}
	case "save":
		storage := NewStorage(config.dbUrl)
		db, err := storage.getDb(config.dbName, config.dbUsername, config.dbPassword)
		if err != nil {
			pretty.Println(err)
			os.Exit(0)
		}
		element := ChoiceElement{Lat: *latPtr, Lng: *lngPtr, Time: time.Now()}
		err = storage.insertChoice(db, element)
		if err != nil {
			pretty.Println(err)
			os.Exit(0)
		}
		pretty.Println("store element: ", element, "to mongodb")
	case "load":
		storage := NewStorage(config.dbUrl)
		db, err := storage.getDb(config.dbName, config.dbUsername, config.dbPassword)
		if err != nil {
			pretty.Println(err)
			os.Exit(0)
		}
		discoverInfo, err := storage.findDiscoverInfo(db, *latPtr, *lngPtr, 1000)
		pretty.Println("discoverInfo: ", discoverInfo)
		choices := storage.findChoiceListByLocation(db, *latPtr, *lngPtr, 1000)
		pretty.Println("choices: ", choices)
	case "web":
		runWebServer()
	case "api":
		runApiServer()
	}

	os.Exit(0)
}
