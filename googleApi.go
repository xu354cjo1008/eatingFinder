/****************************************************************************
 * This file is handler of Geocode api processing.                          *
 *                                                                          *
 ****************************************************************************/
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"googlemaps.github.io/maps"
)

const GOOGLE_API_KEY string = "AIzaSyDJXVVPUtvmRDcBN4nTPNVAI26cUzOaztw"
const GOOGLE_GEOCODE_URL string = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s&language=%s"

/**
 * Interface of geocode api
 */
type googleMapGeocode interface {
	request() error
	getCity() (string, error)
}

/**
 * Class to handle direct access google geocode api
 */
type directGeo struct {
	requestContext string
	response       map[string]interface{}
}

/**
 * Request google geocode api and store response to struct
 */
func (geo *directGeo) request() error {

	if geo.requestContext == "" {
		return errors.New("Invalid request context")
	}

	resp, err := http.Get(geo.requestContext)
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
func newDirectGeo(lat float64, lng float64, language string) *directGeo {
	reqUrl := fmt.Sprintf(GOOGLE_GEOCODE_URL, lat, lng, GOOGLE_API_KEY, language)

	directGeo := directGeo{
		requestContext: reqUrl,
	}

	return &directGeo
}

/**
 * Class to handle google geocode api provide by https://github.com/googlemaps/google-maps-services-go
 */
type mapGeo struct {
	client         *maps.Client
	requestContext *maps.GeocodingRequest
	response       []maps.GeocodingResult
}

/**
 * Request google google map geocode api and store response to struct
 */
func (m *mapGeo) request() error {

	if m.requestContext == nil {
		return errors.New("invalid google map geocoding request")
	}

	resp, err := m.client.ReverseGeocode(context.Background(), m.requestContext)
	if err != nil {
		return err
	}

	m.response = resp

	return nil
}

/**
 * Parse googlemap.map geocode return result
 */
func (mapGeo *mapGeo) getCity() (string, error) {

	var components []maps.AddressComponent
	if len(mapGeo.response) == 0 {
		return "", errors.New("Can not get related address of that location")
	}
	for _, geocodingRes := range mapGeo.response {
		if geocodingRes.Types[0] == "street_address" {
			components = geocodingRes.AddressComponents
			break
		}
	}

	if components == nil || len(components) == 0 {
		return "", errors.New("Can not get related address information")
	}

	// Try to get city name with highest administrative area level
	for _, component := range components {
		if component.Types[0] == "administrative_area_level_1" {
			return component.LongName, nil
		}
	}

	for _, component := range components {
		if component.Types[0] == "administrative_area_level_2" {
			return component.LongName, nil
		}
	}

	for _, component := range components {
		if component.Types[0] == "administrative_area_level_3" {
			return component.LongName, nil
		}
	}

	return "", errors.New("Can not find related city name")
}

/**
 * Contructure of https://github.com/googlemaps/google-maps-services-go geocode class
 */
func newMapGeo(lat float64, lng float64, language string) *mapGeo {
	client, err := maps.NewClient(maps.WithAPIKey(GOOGLE_API_KEY))

	if err != nil {
		return nil
	}
	req := &maps.GeocodingRequest{
		LatLng:   &maps.LatLng{Lat: lat, Lng: lng},
		Language: language,
	}

	mapGeo := mapGeo{
		client:         client,
		requestContext: req,
	}

	return &mapGeo
}

/**
 * @name GetCityByLatlng
 * @brief Get city name by latitude and longtitude
 * @param lat Latitude
 * @param lng Longtitude
 * @return string City name
 * @return error Error description, this will be nil if no error occurs
 */
func GetCityByLatlng(lat float64, lng float64, language string) (string, error) {

	var geo googleMapGeocode

	geo = newMapGeo(lat, lng, language)

	err := geo.request()
	if err != nil {
		return "", err
	}

	var city string = ""
	city, err = geo.getCity()
	if err != nil {
		return "", err
	}

	return city, nil
}
