package nearPlace

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"reflect"

	"github.com/xu354cjo1008/eatingFinder/httpHandler"
)

const googleLib string = "google_lib"
const googleDir string = "google_dir"
const DEF_OUTPUT string = "/json?"
const DEF_TYPE string = "food"
const DEF_LANG string = "en"
const DEF_RANK string = "prominence"
const NSP string = "Not Support Now"
const API_KEY string = "AIzaSyAw-XcK-bqStWykfu3n-kbAeTwJhruRCBc"

type GoogleBase struct {
	handler googleMethod
	source  string
}
type googleMethod interface {
	requireTo() error
	requireBy(float64, float64, uint, string) error
	parsing() ([][]map[string]interface{}, error)
}

var priority = map[string]int{
	googleLib: 1,
	//googleDir: 2,
}
var EMP = map[string]interface{}{
	"uint":    0,
	"int":     0,
	"float64": 0.0,
	"string":  "",
	"bool":    false,
	"error":   nil,
}

type placeInfo struct {
	search searchOpt
}

type searchOpt struct {
	nearby nearbySearch
}

func check(err error) {
	if err != nil {
		log.Fatalf("[NearbySearch]fatal error: %s", err)
	}
}

func whichSource() (res string) {
	res = googleLib
	return
}

/**
 * @name GetNearRestaurants
 * @brief Return the restaurants near the latitude and longtitude.
 * @param lat The latitude.
 * @param lng The longtitude.
 * @param rad The radius.
 * @param lan The language.
 * @return res The restaurants by json.
 * @return err Error description, this will be nil if no error occurs.
 */
func (base *GoogleBase) GetNearRestaurants(lat float64, lng float64, rad uint, lan string) (res [][]map[string]interface{}, err error) {
	// Initialize
	switch base.source {
	case googleLib:
		err = base.handler.requireBy(lat, lng, rad, lan)
		check(err)
		err = base.handler.requireTo()
		check(err)
		res, err = base.handler.parsing()
		check(err)
	case googleDir:
		url := placeUrl["nearbySearch"] + DEF_OUTPUT
		config := placeInfo{
			search: searchOpt{
				nearby: nearbySearch{
					Key: paraFormat{
						para:  "key",
						value: API_KEY,
					},
					Location: paraFormat{
						para:  "location",
						value: fmt.Sprintf("%f, %f", lat, lng),
					},
					Radius: paraFormat{
						para:  "radius",
						value: fmt.Sprintf("%d", rad),
					},
					Keyword: paraFormat{
						para:  "keyword",
						value: NSP,
					},
					Language: paraFormat{
						para:  "language",
						value: DEF_LANG,
					},
					Minprice: paraFormat{
						para:  "minprice",
						value: NSP,
					},
					Maxprice: paraFormat{
						para:  "maxprice",
						value: NSP,
					},
					Name: paraFormat{
						para:  "name",
						value: NSP,
					},
					Opennow: paraFormat{
						para:  "opennow",
						value: NSP,
					},
					Rankby: paraFormat{
						para:  "rankby",
						value: DEF_RANK,
					},
					Type: paraFormat{
						para:  "types",
						value: DEF_TYPE,
					},
					Pagetoken: paraFormat{
						para:  "pagetoken",
						value: NSP,
					},
					Zagatselected: paraFormat{
						para:  "zagatselected",
						value: NSP,
					},
				},
			},
		}
		// Construction the url
		v := reflect.ValueOf(config.search.nearby)
		for i := 0; i < v.NumField(); i++ {
			url = addParameter(url, v.Field(i).Interface().(paraFormat))
		}
		fmt.Println(url)
		// Send http request
		resp, err := httpHandler.HttpGet(url)
		check(err)

		if err := json.Unmarshal(resp, &res); err != nil {
			return nil, err
		}
		fmt.Println(res)
	default:
	}
	return
}

/**
 * @name InitPlaceNearbySearch
 * @brief Create the default resource for communitation with NearbySearch of GoogleMap API.
 * @param source The nearby search source.
 * @return res The basic resource of specific source.
 * @return err Error description, this will be nil if no error occurs.
 */
func InitPlaceNearbySearch(source string) (res *GoogleBase, err error) {
	res = new(GoogleBase)
	res.source = source
	switch source {
	case googleLib:
		res.handler, err = Init_gMapNS(API_KEY, EMP["string"].(string), EMP["string"].(string), EMP["float64"].(float64), EMP["float64"].(float64), uint(EMP["uint"].(int)), EMP["string"].(string), EMP["string"].(string), EMP["string"].(string), EMP["string"].(string), EMP["string"].(string), true, EMP["string"].(string), DEF_TYPE, EMP["string"].(string))
	case googleDir:
	default:
		err = errors.New(fmt.Sprintf("Unknow source: \"%S\"", source))
	}
	return
}
