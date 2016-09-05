package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const GOOGLE_GEOCODE_URL string = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s&language=zh-TW"
const GOOGLE_API_KEY string = ""

func httpGet(request string) (string, error){
	resp, err := http.Get(request)
	if err != nil {
		fmt.Println("http.get failed")
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read http response failed")
		return "", err
	}

	fmt.Println(string(body))
	return string(body), nil
}

func GetCityByLatlng(lat float64, lng float64) (string, error) {
	reqUrl := fmt.Sprintf(GOOGLE_GEOCODE_URL, lat, lng, GOOGLE_API_KEY)
	resp, err := httpGet(reqUrl)
	if err != nil {
		fmt.Printf("error: %v", err)
		return "", err
	}

	return resp, nil
}
