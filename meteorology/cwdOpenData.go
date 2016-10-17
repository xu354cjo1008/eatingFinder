/****************************************************************************
 * This file is xml parser for the data from Central Weather Bureau.        *
 *                                                                          *
 ****************************************************************************/
package meteorology

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/xu354cjo1008/eatingFinder/httpHandler"
)

const (
	CENTRAL_WEATHER_BUREAU_URL       string = "http://opendata.cwb.gov.tw/opendataapi?dataid=%s&authorizationkey=%s"
	CENTRAL_WEATHER_BUREAU_DATA_ID_1 string = "F-C0032-002"
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
	EndTime   string    `xml:"endTime"`
	Parameter parameter `xml:"parameter"`
}

type parameter struct {
	XMLName xml.Name `xml:"parameter"`
	Name    string   `xml:"parameterName"`
	Value   int      `xml:"parameterValue"`
}

type cwdMeteo struct {
	apiKey   string
	language string
	logLevel int
	logger   *log.Logger
	dataset  dataset
}

/**
 * @name DataOfLocation
 * @briefFind weather information from related location
 * @param dataset The dataset struct from xml
 * @param location The location we care about
 * @return *location The pointer of location data from xml
 * @return error The Error description, this will be nil if no error occurs
 */
func (meteo *cwdMeteo) dataOfLocation(dataset dataset, location string) (*location, error) {
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

func (meteo *cwdMeteo) getElement(dataOflocation location, name string) (*weatherElement, error) {

	for _, element := range dataOflocation.WeatherElements {
		if element.ElementName == name {
			return &element, nil
		}
	}
	return nil, errors.New("can not find element with related name")
}

func (meteo *cwdMeteo) getInfoByTime(element weatherElement, inTime time.Time) (*dataByTime, error) {

	if meteo.logLevel == 1 {
		meteo.logger.Println("data time is ", inTime.String())
	}

	for index, dataOfTime := range element.Time {
		startTime, err := time.Parse(time.RFC3339, dataOfTime.StartTime)
		if err != nil {
			return nil, err
		}
		endTime, err := time.Parse(time.RFC3339, dataOfTime.EndTime)
		if err != nil {
			return nil, err
		}

		if meteo.logLevel == 1 {
			meteo.logger.Println("element#", index, " start time is ", startTime.String())
			meteo.logger.Println("element#", index, " end time is ", endTime.String())
		}

		if inTime.After(startTime) && inTime.Before(endTime) {
			return &dataOfTime, nil
		}
	}

	return nil, errors.New("can not find data for that time")
}

func (meteo *cwdMeteo) getParameter(location location, time time.Time, name string) (*dataByTime, error) {

	element, err := meteo.getElement(location, name)
	if err != nil {
		return nil, err
	}

	data, err := meteo.getInfoByTime(*element, time)
	if err != nil {
		return nil, err
	}

	return data, nil
}

/**
 * Parsing weather information from Central Weather Bureau
 * The example xml file is in F-C0032-001.xml and F-C0032-002.xml
 */
func (meteo *cwdMeteo) request() (*Weathers, error) {
	reqUrl := fmt.Sprintf(CENTRAL_WEATHER_BUREAU_URL, CENTRAL_WEATHER_BUREAU_DATA_ID_1, meteo.apiKey)
	resp, err := httpHandler.HttpGet(reqUrl)
	if err != nil {
		return nil, err
	}

	v := Weathers{}
	err = xml.Unmarshal(resp, &v)
	if err != nil {
		return nil, err
	}

	if meteo.logLevel == 1 {
		t, _ := json.MarshalIndent(v, "", "  ")
		meteo.logger.Println("request raw data", string(t))
	}

	return &v, nil
}

func (meteo *cwdMeteo) transformCIToEnum(desc string) int {

	CIMap := map[string]int{
		"COMFORTABLE": CI_COMFORTABLE,
		"HOT":         CI_HOT,
	}

	res := 0

	for index, ci := range CIMap {
		if strings.Contains(desc, index) {
			res += ci
		}
	}

	return res
}

func (meteo *cwdMeteo) getWeather(location string, time time.Time) (*Weather, error) {

	weatherData, err := meteo.request()
	if err != nil {
		return nil, err
	}

	dataOfLocation, err := meteo.dataOfLocation(weatherData.DataSet, location)
	if err != nil {
		return nil, err
	}

	wx, err := meteo.getParameter(*dataOfLocation, time, "Wx")
	if err != nil {
		return nil, err
	}
	maxT, err := meteo.getParameter(*dataOfLocation, time, "MaxT")
	if err != nil {
		return nil, err
	}
	minT, err := meteo.getParameter(*dataOfLocation, time, "MinT")
	if err != nil {
		return nil, err
	}
	ci, err := meteo.getParameter(*dataOfLocation, time, "CI")
	if err != nil {
		return nil, err
	}
	pop, err := meteo.getParameter(*dataOfLocation, time, "PoP")
	if err != nil {
		return nil, err
	}

	maxTemp, err := strconv.Atoi(maxT.Parameter.Name)
	minTemp, err := strconv.Atoi(minT.Parameter.Name)
	probOfprecip, err := strconv.Atoi(pop.Parameter.Name)
	weather := Weather{
		weather:      transformWxToEnum(wx.Parameter.Name),
		maxTemp:      maxTemp,
		minTemp:      minTemp,
		comfortIndex: meteo.transformCIToEnum(ci.Parameter.Name),
		pop:          probOfprecip,
	}

	return &weather, nil
}

func newCwdMeteo(apiKey string, language string, logFile io.Writer) *cwdMeteo {

	var loggingLevel int
	if logFile == nil {
		loggingLevel = 0
	} else {
		loggingLevel = 1
	}
	meteo := cwdMeteo{
		apiKey:   apiKey,
		language: language,
		logLevel: loggingLevel,
		logger: log.New(logFile, "CwdMeteo: ",
			log.Ldate|log.Ltime|log.Lshortfile),
	}

	return &meteo
}
