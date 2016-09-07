package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

const GOOGLE_GEOCODE_URL string = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s&language=zh-TW"
const GOOGLE_API_KEY string = ""

func httpGet(request string) ([]byte, error){
	resp, err := http.Get(request)
	if err != nil {
		fmt.Println("http.get failed")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read http response failed")
		return nil, err
	}

	return body, nil
}

func GetCityByLatlng(lat float64, lng float64) (string, error) {
	reqUrl := fmt.Sprintf(GOOGLE_GEOCODE_URL, lat, lng, GOOGLE_API_KEY)
	resp, err := httpGet(reqUrl)
	if err != nil {
		fmt.Printf("error: %v", err)
		return "", err
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(resp, &dat); err != nil {
		panic(err)
    }

	address := dat["results"].([]interface{})[0].(map[string]interface{})

    components := address["address_components"].([]interface{})

	for _,component := range components {
		types := component.(map[string]interface{})["types"].([]interface{})
		if types[0] == "administrative_area_level_1" {
			return component.(map[string]interface{})["long_name"].(string), nil
		}
	}

	for _,component := range components {
		types := component.(map[string]interface{})["types"].([]interface{})
		if types[0] == "administrative_area_level_2" {
			return component.(map[string]interface{})["long_name"].(string), nil
		}
	}

	return "", nil
}
