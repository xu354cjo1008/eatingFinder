/****************************************************************************
 * This file is handler of google direct geocode api processing.            *
 *                                                                          *
 ****************************************************************************/
package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const GOOGLE_GEOCODE_URL string = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s&language=%s"

/**
 * Class to handle direct access google geocode api
 */
type directGeo struct {
	response     map[string]interface{}
	googleApiKey string
	language     string
}

/**
 * Request google geocode api and store response to struct
 */
func (geo *directGeo) request(lat float64, lng float64) error {

	if geo.googleApiKey == "" {
		return errors.New("Invalid google api key")
	}

	reqUrl := fmt.Sprintf(GOOGLE_GEOCODE_URL, lat, lng, geo.googleApiKey, geo.language)

	resp, err := http.Get(reqUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var geoRes map[string]interface{}
	if err := json.Unmarshal(body, &geoRes); err != nil {
		return err
	}

	geo.response = geoRes

	return nil
}

/**
 * Parse geocode api return result
 * The example geocode result is at geocodeReturnExample file
 */
func (geo *directGeo) getCity() (string, error) {

	if geo.response == nil {
		return "", errors.New("Can not get city from Invalid response context")
	}

	result := geo.response["results"].([]interface{})
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
 * Contructure of direct geocode class
 */
func newDirectGeo(googleApiKey string, language string) *directGeo {

	geo := directGeo{
		googleApiKey: googleApiKey,
		language:     language,
	}

	return &geo
}
