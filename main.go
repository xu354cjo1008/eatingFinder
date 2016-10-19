package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"github.com/xu354cjo1008/eatingFinder/geography/geocoding"
	"github.com/xu354cjo1008/eatingFinder/httpHandler"
	"github.com/xu354cjo1008/eatingFinder/meteorology"
)

func meteoUtil(lat float64, lng float64, logFile string, googleApiKey string, cwdApiKey string) error {

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

	geocode := geocoding.NewGeocode(googleApiKey, "en")

	city, err := geocode.GetCityByLatlng(lat, lng)

	if err != nil {
		fmt.Println("error: ", err)
		return err
	}

	pretty.Println(city)

	meteo := meteorology.NewMeteorology(cwdApiKey, "en", file)
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
	var apiHost string
	var apiPort int
	var googleApiKey string
	var cwdApiKey string

	mode := flag.String("mode", "meteo", "utility mode: <meteo|server>")
	latPtr := flag.Float64("lat", 25.057339, "latitude of user position")
	lngPtr := flag.Float64("lng", 121.56086, "longtitude of user position")
	logFilePtr := flag.String("log", "", "log path <path|fg>")

	flag.Parse()

	viper.SetConfigName("app")
	viper.AddConfigPath("config")

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Config file not found...")
	} else {
		apiHost = viper.GetString("development.apiHost")
		apiPort = viper.GetInt("development.apiPort")
		googleApiKey = viper.GetString("development.googleApiKey")
		cwdApiKey = viper.GetString("development.cwdApiKey")
	}

	log.Printf("\nDevelopment Config found:\n server = %s\n"+
		" port = %d\n",
		apiHost,
		apiPort)

	switch *mode {
	case "meteo":
		err = meteoUtil(*latPtr, *lngPtr, *logFilePtr, googleApiKey, cwdApiKey)
		if err != nil {
			pretty.Println(err)
			os.Exit(-1)
		}
	case "server":
		httpHandler.RunServer(apiPort)
	}

	os.Exit(0)
}
