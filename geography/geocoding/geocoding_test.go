/****************************************************************************
 * The unit tester for geocoding go package                          .      *
 *                                                                          *
 ****************************************************************************/
package geocoding

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
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

	var googleApiKey string
	var err error

	viper.SetConfigName("app")
	viper.AddConfigPath("../../config")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...")
	} else {
		googleApiKey = viper.GetString("development.googleApiKey")
	}

	testCases := []locationTestCase{
		locationTestCase{lat: 25.053257, lng: 121.539702, expect: "Taipei City"},
		locationTestCase{lat: 24.744071, lng: 121.763291, expect: "Yilan County"},
		locationTestCase{lat: 23.968259, lng: 121.583281, expect: "Hualien County"},
		locationTestCase{lat: 23.569817, lng: 119.640066, expect: "Penghu County"},
		locationTestCase{lat: 25.129331, lng: 121.739967, expect: "Keelung City"},
	}

	geocode := NewGeocode(googleApiKey, "en")

	for index, testCase := range testCases {
		if res, err := geocode.GetCityByLatlng(testCase.lat, testCase.lng); res != testCase.expect || err != nil {
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
