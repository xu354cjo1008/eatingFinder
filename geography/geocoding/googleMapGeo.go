/****************************************************************************
 * This file is handler of google map geocode api processing.               *
 * Implement based on googlemaps.github.io/maps                             *
 ****************************************************************************/
package geocoding

import (
	"context"
	"errors"

	"googlemaps.github.io/maps"
)

/**
 * Class to handle google geocode api provide by https://github.com/googlemaps/google-maps-services-go
 */
type mapGeo struct {
	client   *maps.Client
	response []maps.GeocodingResult
	language string
}

/**
 * Request google google map geocode api and store response to struct
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
