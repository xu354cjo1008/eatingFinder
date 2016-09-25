package geocoding

import (
	"context"
	"errors"
	"googlemaps.github.io/maps"
)

const GOOGLE_GEOCODE_URL string = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s&language=%s"

type Geocode struct {
	geoHandler   googleMapGeocode
	googleApiKey string
	language     string
}

/**
 * Interface of geocode api
 */
type googleMapGeocode interface {
	request(float64, float64) error
	getCity() (string, error)
}

/**
 * @name GetCityByLatlng
 * @brief Get city name by latitude and longtitude
 * @param lat Latitude
 * @param lng Longtitude
 * @return string City name
 * @return error Error description, this will be nil if no error occurs
 */
func (geo *Geocode) GetCityByLatlng(lat float64, lng float64) (string, error) {

	err := geo.geoHandler.request(lat, lng)
	if err != nil {
		return "", err
	}

	var city string = ""
	city, err = geo.geoHandler.getCity()
	if err != nil {
		return "", err
	}

	return city, nil
}

/**
 * Contructure of https://github.com/googlemaps/google-maps-services-go geocode class
 */
func newMapGeo(googleApiKey string, language string) *mapGeo {

	client, err := maps.NewClient(maps.WithAPIKey(googleApiKey))

	if err != nil {
		return nil
	}

	geo := mapGeo{
		client:   client,
		language: language,
	}

	return &geo
}

/**
 * @name NewGeoCode
 * @brief Create a geocode instance
 * @param googleApiKey google map api key
 * @param language language e.g. en, zh-TW
 * @return Geocode instance
 */
func NewGeocode(googleApiKey string, language string) *Geocode {

	geo := Geocode{
		geoHandler:   newMapGeo(googleApiKey, language),
		googleApiKey: googleApiKey,
		language:     language,
	}

	return &geo
}

/**
 * Parse googlemap.map geocode return result
 */
func (geo *mapGeo) getCity() (string, error) {

	var components []maps.AddressComponent
	if len(geo.response) == 0 {
		return "", errors.New("Can not get related address of that location")
	}

	components = geo.response[0].AddressComponents

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
 * Class to handle direct access google geocode api
 */
type directGeo struct {
	response     map[string]interface{}
	googleApiKey string
	language     string
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

/**
 * Class to handle google geocode api provide by https://github.com/googlemaps/google-maps-services-go
 */
type mapGeo struct {
	client   *maps.Client
	response []maps.GeocodingResult
	language string
}

/**
 * Request google map geocode api and store response to struct
 */
func (geo *mapGeo) request(lat float64, lng float64) error {
	req := &maps.GeocodingRequest{
		LatLng:   &maps.LatLng{Lat: lat, Lng: lng},
		Language: geo.language,
	}

	resp, err := geo.client.ReverseGeocode(context.Background(), req)
	if err != nil {
		return err
	}

	geo.response = resp

	return nil
}
