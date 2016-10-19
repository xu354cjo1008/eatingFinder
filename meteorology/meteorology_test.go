/****************************************************************************
 * The unit tester for meteorology go package                          .    *
 *                                                                          *
 ****************************************************************************/
package meteorology

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

type weatherTestCase struct {
	location string
	expect   int
}

/**
 * Test job for weather information parser
 */
func TestWeatherApi(t *testing.T) {

	var cwdApiKey string
	var err error

	viper.SetConfigName("app")
	viper.AddConfigPath("../config")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...")
	} else {
		cwdApiKey = viper.GetString("development.cwdApiKey")
	}

	testCases := []weatherTestCase{
		weatherTestCase{location: "Taipei City", expect: WX_MOSTLY + WX_CLEAR},
		weatherTestCase{location: "New Taipei City", expect: WX_MOSTLY + WX_CLEAR},
		weatherTestCase{location: "Taoyuan City", expect: WX_MOSTLY + WX_CLEAR},
	}

	for index, testCase := range testCases {
		//		dataOfLocation, err := DataOfLocation(weatherData.DataSet, testCase.location)
		//		if dataOfLocation.WeatherElements[0].Time[0].Parameter.Name != testCase.expect || err != nil {
		meteo := NewMeteorology(cwdApiKey, "en", nil)
		data, err := meteo.GetWeather(testCase.location)
		if err != nil || data.weather != testCase.expect {
			t.Error(
				"#", index,
				"location", testCase.location,
				"Expected", testCase.expect,
				"Got", data.weather,
				"Failed",
			)
		} else {
			t.Log(
				"#", index,
				"location", testCase.location,
				"Expected", testCase.expect,
				"Got", data.weather,
				"Pass",
			)
		}
	}
}
