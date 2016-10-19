package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
	}

	log.Printf("\nDevelopment Config found:\n default server port = %d\n"+
		" api host = %s\n"+
		" api port = %d\n",
		config.defaultPort,
		config.apiHost,
		config.apiPort)

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

/**
 * This is the main just for test
 * We need to write another unit test program to do this
 */
func main() {

	var err error

	mode := flag.String("mode", "meteo", "utility mode: <meteo|web|api>")
	latPtr := flag.Float64("lat", 25.057339, "latitude of user position")
	lngPtr := flag.Float64("lng", 121.56086, "longtitude of user position")
	logFilePtr := flag.String("log", "", "log path <path|fg>")

	flag.Parse()

	err = configure()

	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	switch *mode {
	case "meteo":
		err = meteoUtil(*latPtr, *lngPtr, *logFilePtr)
		if err != nil {
			pretty.Println(err)
			os.Exit(-1)
		}
	case "web":
		runWebServer()
	case "api":
		runApiServer()
	}

	os.Exit(0)
}
