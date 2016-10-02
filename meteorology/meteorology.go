/****************************************************************************
 * This file is meteorology core function and structure.                    *
 *                                                                          *
 ****************************************************************************/
package meteorology

import "time"

/**
 * enum bitmap of weather description
 * e.g. 1000 0010 means party cloudy
 */
const (
	WX_CLEAR = 1 << iota
	WX_CLOUDY
	WX_FOG
	WX_RAIN
	WX_SHOWERS
	WX_THUNDERSTORMS
	WX_THUNDERSHOWERS
	WX_PARTLY
	WX_MOSTLY
	WX_OCCASIONAL
	WX_LOCAL
	WX_AFTERNOON
)

const (
	CI_COMFORTABLE = 1 << iota
	CI_HOT
)

type meteorology interface {
	getWeather(string, time.Time) (*Weather, error)
}

type Meteorology struct {
	meteoHandler meteorology
	apiKey       string
	language     string
}

type Weather struct {
	weather      int
	maxTemp      int
	minTemp      int
	comfortIndex int
	pop          int
}

func (meteo *Meteorology) GetWeather(location string) (*Weather, error) {
	t := time.Now()
	data, err := meteo.meteoHandler.getWeather(location, t)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewMeteorology(apiKey string, language string) *Meteorology {

	meteo := Meteorology{
		meteoHandler: newCwdMeteo(apiKey, language),
		apiKey:       apiKey,
		language:     language,
	}

	return &meteo
}
