package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/creack/goproxy"
	"github.com/creack/goproxy/registry"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/xu354cjo1008/eatingFinder/geography/geocoding"
)

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Home")
}

func apiGeocodeHandler(rw http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(rw, "Api Geocode Handler")

	vars := mux.Vars(r)

	geocode := geocoding.NewGeocode(config.googleApiKey, "en")

	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lng, _ := strconv.ParseFloat(vars["lng"], 64)

	city, err := geocode.GetCityByLatlng(lat, lng)

	if err != nil {
		log.Println("error: ", err)
	}

	fmt.Fprintln(rw, city)
}

func runApiServer() {

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/getCity/{lat}/{lng}", apiGeocodeHandler)

	n := negroni.Classic()
	n.UseHandler(r)

	log.Println("Starting api server on " + ":" + strconv.Itoa(config.defaultPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.defaultPort), n))
}

func runWebServer() {

	apiRegistry := registry.DefaultRegistry{
		"api": {
			"v1": {
				config.apiHost + ":" + strconv.Itoa(config.apiPort),
			},
		},
	}

	proxy := goproxy.NewMultipleHostReverseProxy(apiRegistry)

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.PathPrefix("/api").Handler(proxy)

	n := negroni.Classic()
	n.UseHandler(r)

	log.Println("Starting web server on " + ":" + strconv.Itoa(config.defaultPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.defaultPort), n))
}
