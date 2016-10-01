/****************************************************************************
 * This file is meteorology core function and structure.                    *
 *                                                                          *
 ****************************************************************************/
package meteorology

import "time"

type meteorology interface {
	getWeather(string, time.Time) (*Weather, error)
}

type Meteorology struct {
	meteoHandler meteorology
	apiKey       string
	language     string
}

type Weather struct {
	weather      string
	maxTemp      int
	minTemp      int
	comfortIndex string
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
