package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"errors"
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
	for _,data := range dataset.Locations {
		if data.LocationName == location {
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

	v := parseWeatherXml("F-C0032-001.xml")

	dataOfLocation, _ := dataOfLocation(v.DataSet, "臺北市")
	fmt.Println(dataOfLocation.LocationName)
	fmt.Println(dataOfLocation.WeatherElements[0].Time[0].StartTime)
	fmt.Println(dataOfLocation.WeatherElements[0].Time[0].Parameter.Name)
}
