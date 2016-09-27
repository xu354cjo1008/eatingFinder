/****************************************************************************
 * The unit tester for meteorology go package                          .    *
 *                                                                          *
 ****************************************************************************/
package meteorology

import "testing"

type weatherTestCase struct {
	location string
	expect   string
}

/**
 * Test job for weather information parser
 */
func TestWeatherApi(t *testing.T) {
	testCases := []weatherTestCase{
		weatherTestCase{location: "Taipei City", expect: "MOSTLY CLOUDY WITH SHOWERS OR THUNDERSTORMS"},
		weatherTestCase{location: "New Taipei City", expect: "MOSTLY CLOUDY WITH SHOWERS OR THUNDERSTORMS"},
		weatherTestCase{location: "Taoyuan City", expect: "CLOUDY WITH SHOWERS OR THUNDERSTORMS"},
	}
	weatherData, err := ParseWeatherXml()
	if err != nil {
		t.Error(
			"Failed to parse weather xml",
		)
		return
	}
	for index, testCase := range testCases {
		dataOfLocation, err := DataOfLocation(weatherData.DataSet, testCase.location)
		if dataOfLocation.WeatherElements[0].Time[0].Parameter.Name != testCase.expect || err != nil {
			t.Error(
				"#", index,
				"location", testCase.location,
				"Expected", testCase.expect,
				"Got", dataOfLocation.WeatherElements[0].Time[0].Parameter.Name,
				"Failed",
			)
		} else {
			t.Log(
				"#", index,
				"location", testCase.location,
				"Expected", testCase.expect,
				"Got", dataOfLocation.WeatherElements[0].Time[0].Parameter.Name,
				"Pass",
			)
		}
	}
}
