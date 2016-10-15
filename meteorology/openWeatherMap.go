/****************************************************************************
 * This file is handler for openWeatherMap                                  *
 * Related API information is at http://openweathermap.org/                 *
 ****************************************************************************/
package meteorology

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kr/pretty"
	"github.com/xu354cjo1008/eatingFinder/httpHandler"
)

const (
	OPEN_WEATHER_MAP_API      string = "http://api.openweathermap.org/data/2.5/%s?q=%s,%s&APPID=%s"
	OPEN_WEATHER_MAP_WEATHER  string = "weather"
	OPEN_WEATHER_MAP_FORECAST string = "forecast"
)

type owmMeteo struct {
	apiKey   string
	language string
	response map[string]interface{}
}

func (meteo *owmMeteo) request(city string, country string, reqType string) error {

	if meteo.apiKey == "" {
		return errors.New("Invalid openWeatherMap api key")
	}
	if city == "" {
		return errors.New("Invalid city name")
	}
	if country == "" {
		return errors.New("Invalid country name")
	}

	var requestT string

	switch reqType {
	case "weather":
		requestT = OPEN_WEATHER_MAP_WEATHER
	case "forecast":
		requestT = OPEN_WEATHER_MAP_FORECAST
	}

	reqUrl := fmt.Sprintf(OPEN_WEATHER_MAP_API, requestT, city, country, meteo.apiKey)

	resp, err := httpHandler.HttpGet(reqUrl)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &(meteo.response)); err != nil {
		return err
	}

	return nil
}

func (meteo *owmMeteo) fakeRequest(city string, country string, reqType string) error {

	if meteo.apiKey == "" {
		return errors.New("Invalid openWeatherMap api key")
	}
	if city == "" {
		return errors.New("Invalid city name")
	}
	if country == "" {
		return errors.New("Invalid country name")
	}

	aa, _ := ioutil.ReadFile("testCase/openWeatherMapTaipeiForecast.json")
	var testaa map[string]interface{}
	if err := json.Unmarshal(aa, &testaa); err != nil {
		return err
	}
	pretty.Println(testaa)

	resp, err := ioutil.ReadFile("testCase/openWeatherMapTaipeiWeather.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return err
	}

	if err := json.Unmarshal(resp, &(meteo.response)); err != nil {
		return err
	}
	pretty.Println(meteo.response)

	return nil
}

func (meteo *owmMeteo) getElement(name string) (map[string]interface{}, error) {

	if meteo.response == nil && len(meteo.response) == 0 {
		return nil, errors.New("empty data in owmMeteo response")
	}

	element := meteo.response[name]

	if element == nil {
		return nil, errors.New("empty data in the request element")
	}

	switch name {
	case "weather":
		return element.([]interface{})[0].(map[string]interface{}), nil
	default:
		return element.(map[string]interface{}), nil
	}
}

func (meteo *owmMeteo) getParameter(element map[string]interface{}, name string) (interface{}, error) {

	if element == nil {
		return nil, errors.New("empty data in request element")
	}

	parameter := element[name]
	if parameter == nil {
		return nil, errors.New("empty data in the request parameter")
	}

	return parameter, nil
}

func (meteo *owmMeteo) tempKToCel(temp float64) float64 {
	return temp - 273.15
}

func (meteo *owmMeteo) getWeather(location string, time time.Time) (*Weather, error) {

	err := meteo.fakeRequest("Taipei", "Taiwan", "Weather")
	//	err := meteo.request("Taipei", "Taiwan", "Weather")
	if err != nil {
		return nil, err
	}

	var element map[string]interface{}
	var parameter interface{}
	weather := Weather{}

	element, _ = meteo.getElement("weather")
	pretty.Println(element)
	parameter, _ = meteo.getParameter(element, "description")
	pretty.Println("weather: ", parameter)

	element, _ = meteo.getElement("main")
	parameter, _ = meteo.getParameter(element, "temp_max")
	weather.maxTemp = int(meteo.tempKToCel(parameter.(float64)))
	parameter, _ = meteo.getParameter(element, "temp_min")
	weather.minTemp = int(meteo.tempKToCel(parameter.(float64)))

	return &weather, nil
}

func newOwmMeteo(apiKey string, language string) *owmMeteo {

	meteo := owmMeteo{
		apiKey: apiKey,
	}

	return &meteo
}
