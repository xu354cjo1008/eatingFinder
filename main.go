package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	flag.Parse()

	geocode := geocoding.NewGeocode("AIzaSyDJXVVPUtvmRDcBN4nTPNVAI26cUzOaztw", "en")

	city, err := geocode.GetCityByLatlng(*latPtr, *lntPtr)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}

	pretty.Println(city)

	meteo := meteorology.NewMeteorology("CWB-2FC70596-59B4-4CC5-98E5-BCC6490E30DD", "en")
	data, err := meteo.GetWeather(city)
	if err != nil {
		log.Println("error: ", err)
		os.Exit(-1)
	}

	pretty.Println(data)

	os.Exit(0)
}
