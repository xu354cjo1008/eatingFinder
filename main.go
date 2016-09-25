package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xu354cjo1008/weatherGo/googleApi/geocoding"
	"github.com/xu354cjo1008/weatherGo/meteorological"
)

/**
 * This is the main just for test
 * We need to write another unit test program to do this
 */
func main() {

	latPtr := flag.Float64("lat", 25.057339, "latitude of user position")
	lntPtr := flag.Float64("lnt", 121.56086, "longtitude of user position")

	flag.Parse()

	geocode := NewGeocode("AIzaSyDJXVVPUtvmRDcBN4nTPNVAI26cUzOaztw", "en")

	city, err := geocode.GetCityByLatlng(*latPtr, *lntPtr)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}

	fmt.Println(city)

	v, err := ParseWeatherXml()
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}

	dataOfLocation, err := DataOfLocation(v.DataSet, city)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(-1)
	}

	fmt.Println(dataOfLocation.LocationName)
	fmt.Println(dataOfLocation.WeatherElements[0].Time[0].StartTime)
	fmt.Println(dataOfLocation.WeatherElements[0].Time[0].Parameter.Name)

	os.Exit(0)
}
