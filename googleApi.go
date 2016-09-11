/****************************************************************************
 * This file is handler of Geocode api processing.                          *
 *                                                                          *
 ****************************************************************************/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const GOOGLE_GEOCODE_URL string = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s&language=zh-TW"
const GOOGLE_API_KEY string = ""

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

/*
* Parse geocode api return result
* The example geocode result is at geocodeReturnExample file
 */
func parseGeocodeResult(geoRes map[string]interface{}) (string, error) {

	result := geoRes["results"].([]interface{})
	if len(result) == 0 {
		return "", errors.New("Get zero location result")
	}

	address := result[0].(map[string]interface{})
	if len(address) == 0 {
		return "", errors.New("Can not get related address of that location")
	}

	components := address["address_components"].([]interface{})
	if len(components) == 0 {
		return "", errors.New("Can not get related data components of that address")
	}

	// Try to get city name with highest administrative area level
	for _, component := range components {
		types := component.(map[string]interface{})["types"].([]interface{})
		if types[0] == "administrative_area_level_1" {
			return component.(map[string]interface{})["long_name"].(string), nil
		}
	}

	for _, component := range components {
		types := component.(map[string]interface{})["types"].([]interface{})
		if types[0] == "administrative_area_level_2" {
			return component.(map[string]interface{})["long_name"].(string), nil
		}
	}

	for _, component := range components {
		types := component.(map[string]interface{})["types"].([]interface{})
		if types[0] == "administrative_area_level_3" {
			return component.(map[string]interface{})["long_name"].(string), nil
		}
	}

	return "", errors.New("Can not find related city name")

}

/**
* @name GetCityByLatlng
* @brief Get city name by latitude and longtitude
* @param lat Latitude
* @param lng Longtitude
* @return string City name
* @return error Error description, this will be nil if no error occurs
 */
func GetCityByLatlng(lat float64, lng float64) (string, error) {

	reqUrl := fmt.Sprintf(GOOGLE_GEOCODE_URL, lat, lng, GOOGLE_API_KEY)
	resp, err := httpGet(reqUrl)
	if err != nil {
		fmt.Printf("error: %v", err)
		return "", err
	}

	var geoRes map[string]interface{}
	if err := json.Unmarshal(resp, &geoRes); err != nil {
		panic(err)
	}

	var city string = ""
	city, err = parseGeocodeResult(geoRes)
	if err != nil {
		return "", errors.New("Can not find related city name")
	} else {
		return city, nil
	}
}
