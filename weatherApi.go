/****************************************************************************
 * This file is xml parser for the data from Central Weather Bureau.        *
 *                                                                          *
 ****************************************************************************/
package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	CENTRAL_WEATHER_BUREAU_URL       string = "http://opendata.cwb.gov.tw/opendataapi?dataid=%s&authorizationkey=%s"
	CENTRAL_WEATHER_BUREAU_DATA_ID_1 string = "F-C0032-002"
	CENTRAL_WEATHER_BUREAU_KEY       string = "CWB-2FC70596-59B4-4CC5-98E5-BCC6490E30DD"
)

/**
 * The xml structure of weather information from Central Weather Bureau
 * The example xml file is in F-C0032-001.xml and F-C0032-002.xml
 */
type Weathers struct {
	XMLName xml.Name `xml:"cwbopendata"`
	DataSet dataset  `xml:"dataset"`
}

type dataset struct {
	XMLName   xml.Name   `xml:"dataset"`
	Locations []location `xml:"location"`
}

type location struct {
	XMLName         xml.Name         `xml:"location"`
	LocationName    string           `xml:"locationName"`
	WeatherElements []weatherElement `xml:"weatherElement"`
}

type weatherElement struct {
	XMLName     xml.Name     `xml:"weatherElement"`
	ElementName string       `xml:"elementName"`
	Time        []dataByTime `xml:"time"`
}

type dataByTime struct {
	XMLName   xml.Name  `xml:"time"`
	StartTime string    `xml:"startTime"`
	EndTime   string    `xml:"EndTime"`
	Parameter parameter `xml:"parameter"`
}

type parameter struct {
	XMLName xml.Name `xml:"parameter"`
	Name    string   `xml:"parameterName"`
	Value   int      `xml:"parameterValue"`
}

/**
 * Same as httpGet function in googleApi.go
 * Maybe we can manage http function in one package
 */
func httpGet(request string) ([]byte, error) {
	resp, err := http.Get(request)
	if err != nil {
		fmt.Println("http.get failed")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read http response failed")
		return nil, err
	}

	return body, nil
}

/**
 * @name DataOfLocation
 * @briefFind weather information from related location
 * @param dataset The dataset struct from xml
 * @param location The location we care about
 * @return *location The pointer of location data from xml
 * @return error The Error description, this will be nil if no error occurs
 */
func DataOfLocation(dataset dataset, location string) (*location, error) {
	if location == "" {
		return nil, errors.New("invalid location")
	}

	var buffer bytes.Buffer

	// This is a workround to match 台 and 臺
	if string([]rune(location)[0]) == "台" {
		buffer.WriteRune('臺')
		buffer.WriteString(strings.Split(location, "台")[1])
	} else {
		buffer.WriteString(location)
	}

	for _, data := range dataset.Locations {
		if strings.ToLower(data.LocationName) == strings.ToLower(buffer.String()) {
			return &data, nil
		}
	}

	return nil, errors.New("can not find data for the location")
}

/**
 * Parsing weather information from Central Weather Bureau
 * The example xml file is in F-C0032-001.xml and F-C0032-002.xml
 */
func ParseWeatherXml() (*Weathers, error) {

	reqUrl := fmt.Sprintf(CENTRAL_WEATHER_BUREAU_URL, CENTRAL_WEATHER_BUREAU_DATA_ID_1, CENTRAL_WEATHER_BUREAU_KEY)
	resp, err := httpGet(reqUrl)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}

	v := Weathers{}
	err = xml.Unmarshal(resp, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	return &v, nil
}

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
