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

/**
 * Test job for geocode api function
 * Input: latitude, longtitude
 * Output: expected city name
 * Add test case into testCases array if needed
 */
func TestLocation(t *testing.T) {
	testCases := []locationTestCase{
		locationTestCase{lat: 25.053257, lng: 121.539702, expect: "台北市"},
		locationTestCase{lat: 24.744071, lng: 121.763291, expect: "宜蘭縣"},
		locationTestCase{lat: 23.968259, lng: 121.583281, expect: "花蓮縣"},
		locationTestCase{lat: 23.569817, lng: 119.640066, expect: "澎湖縣"},
	}
	for index, testCase := range testCases {
		if res, err := GetCityByLatlng(testCase.lat, testCase.lng); res != testCase.expect || err != nil {
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
