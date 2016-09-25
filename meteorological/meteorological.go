/****************************************************************************
 * This file is xml parser for the data from Central Weather Bureau.        *
 *                                                                          *
 ****************************************************************************/
package meteorological

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/xu354cjo1008/weatherGo/httpHandler"
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
	resp, err := httpHandler.HttpGet(reqUrl)
	if err != nil {
		return nil, err
	}

	v := Weathers{}
	err = xml.Unmarshal(resp, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
