package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/xu354cjo1008/eatingFinder/geography/geocoding"
	"github.com/xu354cjo1008/eatingFinder/meteorology"
)

/**
 * This is the main just for test
 * We need to write another unit test program to do this
 */
func main() {

	latPtr := flag.Float64("lat", 25.057339, "latitude of user position")
	lntPtr := flag.Float64("lnt", 121.56086, "longtitude of user position")
	logFile := flag.String("log", "", "log path")

	flag.Parse()

	var file io.Writer = nil
	var err error

	if logFile != nil && *logFile != "" {
		if strings.Compare(*logFile, "fg") == 0 {
			file = os.Stdout
		} else {
			file, err = os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalln("Failed to open log file :", err)
			}
		}
	}

	geocode := geocoding.NewGeocode("AIzaSyDJXVVPUtvmRDcBN4nTPNVAI26cUzOaztw", "en")

	city, err := geocode.GetCityByLatlng(*latPtr, *lntPtr)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}

	pretty.Println(city)

	meteo := meteorology.NewMeteorology("CWB-2FC70596-59B4-4CC5-98E5-BCC6490E30DD", "en", file)
	data, err := meteo.GetWeather(city)
	if err != nil {
		log.Println("error: ", err)
		os.Exit(-1)
	}

	pretty.Println(data)

	os.Exit(0)
}
