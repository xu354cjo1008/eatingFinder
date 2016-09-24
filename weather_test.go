/****************************************************************************
 * The unit tester for weather go package                          .        *
 *                                                                          *
 ****************************************************************************/
package main

import (
	"testing"
)

type locationTestCase struct {
	lat    float64
	lng    float64
	expect string
}

type weatherTestCase struct {
	location string
	expect   string
}

/**
 * Test job for geocode api function
 * Input: latitude, longtitude
 * Output: expected city name
 * Add test case into testCases array if needed
 */
func TestLocation(t *testing.T) {
	testCases := []locationTestCase{
		locationTestCase{lat: 25.053257, lng: 121.539702, expect: "Taipei City"},
		locationTestCase{lat: 24.744071, lng: 121.763291, expect: "Yilan County"},
		locationTestCase{lat: 23.968259, lng: 121.583281, expect: "Hualien County"},
		locationTestCase{lat: 23.569817, lng: 119.640066, expect: "Penghu County"},
		locationTestCase{lat: 25.129331, lng: 121.739967, expect: "Keelung City"},
	}
	for index, testCase := range testCases {
		if res, err := GetCityByLatlng(testCase.lat, testCase.lng, "en"); res != testCase.expect || err != nil {
			t.Error(
				"#", index,
				"For latitude", testCase.lat,
				"longtitude", testCase.lng,
				"Expected", testCase.expect,
				"Got", res,
				"Failed",
			)
		} else {
			t.Log(
				"#", index,
				"For latitude", testCase.lat,
				"longtitude", testCase.lng,
				"Expected", testCase.expect,
				"Got", res,
				"Pass",
			)
		}
	}
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
