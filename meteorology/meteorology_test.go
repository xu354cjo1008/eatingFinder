/****************************************************************************
 * The unit tester for meteorology go package                          .    *
 *                                                                          *
 ****************************************************************************/
package meteorology

import "testing"

type weatherTestCase struct {
	location string
	expect   int
}

/**
 * Test job for weather information parser
 */
func TestWeatherApi(t *testing.T) {

	testCases := []weatherTestCase{
		weatherTestCase{location: "Taipei City", expect: WX_MOSTLY + WX_CLEAR},
		weatherTestCase{location: "New Taipei City", expect: WX_MOSTLY + WX_CLEAR},
		weatherTestCase{location: "Taoyuan City", expect: WX_MOSTLY + WX_CLEAR},
	}

	for index, testCase := range testCases {
		//		dataOfLocation, err := DataOfLocation(weatherData.DataSet, testCase.location)
		//		if dataOfLocation.WeatherElements[0].Time[0].Parameter.Name != testCase.expect || err != nil {
		meteo := NewMeteorology("CWB-2FC70596-59B4-4CC5-98E5-BCC6490E30DD", "en")
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
