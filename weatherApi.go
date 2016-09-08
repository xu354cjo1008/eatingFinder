package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"errors"
	"bytes"
	"strings"
	"flag"
)

type Recurlyweathers struct {
	XMLName		xml.Name	`xml:"cwbopendata"`
	DataSet		dataset		`xml:"dataset"`
}

type dataset struct {
	XMLName		xml.Name	`xml:"dataset"`
	Locations	[]location	`xml:"location"`
}

type location struct {
	XMLName         xml.Name			`xml:"location"`
	LocationName    string				`xml:"locationName"`
	WeatherElements	[]weatherElement	`xml:"weatherElement"`
}

type weatherElement struct {
	XMLName			xml.Name		`xml:"weatherElement"`
	ElementName		string			`xml:"elementName"`
	Time			[]dataByTime	`xml:"time"`
}

type dataByTime struct {
	XMLName			xml.Name		`xml:"time"`
	StartTime		string			`xml:"startTime"`
	EndTime			string			`xml:"EndTime"`
	Parameter		parameter		`xml:"parameter"`
}

type parameter struct {
	XMLName			xml.Name		`xml:"parameter"`
	Name			string			`xml:"parameterName"`
	Value			int				`xml:"parameterValue"`
}

func dataOfLocation(dataset dataset, location string) (*location, error) {
	if location == "" {
		return nil, errors.New("invalid location")
	}

	var buffer bytes.Buffer

	if string([]rune(location)[0]) == "台" {
		buffer.WriteRune('臺')
		buffer.WriteString(strings.Split(location, "台")[1])
	} else {
		buffer.WriteString(location)
	}

	for _,data := range dataset.Locations {
		if data.LocationName == buffer.String() {
			return &data, nil
		}
	}

	return nil, errors.New("can not find data for the location")
}

func parseWeatherXml(filename string) *Recurlyweathers {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}
	v := Recurlyweathers{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}
	return &v
}

func main() {

	latPtr := flag.Float64("lat", 25.057339, "latitude of user position")
	lntPtr := flag.Float64("lnt", 121.56086, "longtitude of user position")

	flag.Parse()

	city, err := GetCityByLatlng(*latPtr, *lntPtr)

	if err != nil {
		 fmt.Println(err)
		 return
	}

	fmt.Println(city)

	v := parseWeatherXml("F-C0032-001.xml")

	dataOfLocation, err := dataOfLocation(v.DataSet, city)

	if err != nil {
		 fmt.Println(err)
		 return
	}

	fmt.Println(dataOfLocation.LocationName)
	fmt.Println(dataOfLocation.WeatherElements[0].Time[0].StartTime)
	fmt.Println(dataOfLocation.WeatherElements[0].Time[0].Parameter.Name)
}
