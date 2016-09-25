/****************************************************************************
 * This file is handler of Geocode api processing.                          *
 *                                                                          *
 ****************************************************************************/
package geocoding

/**
 * Interface of geocode api
 */
type googleMapGeocode interface {
	request(float64, float64) error
	getCity() (string, error)
}

type Geocode struct {
	geoHandler   googleMapGeocode
	googleApiKey string
	language     string
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
