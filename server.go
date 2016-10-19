package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

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

type Registry map[string][]string

func extractNameVersion(target *url.URL) (name, version string, err error) {
	path := target.Path
	// Trim the leading `/`
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	// Explode on `/` and make sure we have at least
	// 2 elements (service name and version)
	tmp := strings.Split(path, "/")
	if len(tmp) < 2 {
		return "", "", fmt.Errorf("Invalid path")
	}
	name, version = tmp[0], tmp[1]
	// Rewrite the request's path without the prefix.
	target.Path = "/" + strings.Join(tmp[2:], "/")
	return name, version, nil
}

func newMultipleHostReverseProxy(reg Registry) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		name, version, err := extractNameVersion(req.URL)
		if err != nil {
			log.Print(err)
			return
		}
		req.URL.Scheme = "http"
		req.URL.Host = name + "/" + version
		log.Println("redrect to http://" + req.URL.Host + req.URL.Path)
	}
	return &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: func(network, addr string) (net.Conn, error) {
				// Trim the `:80` added by Scheme http.
				addr = strings.Split(addr, ":")[0]
				endpoints := reg[addr]
				if len(endpoints) == 0 {
					return nil, fmt.Errorf("Service/Version not found")
				}
				return net.Dial(network, endpoints[rand.Int()%len(endpoints)])
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func runApiServer() {

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/getCity/{lat}/{lng}", apiGeocodeHandler)

	n := negroni.Classic()
	n.UseHandler(r)

	log.Println("Starting api server on " + ":" + strconv.Itoa(config.apiPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.apiPort), n))
}

func runWebServer() {

	proxy := newMultipleHostReverseProxy(Registry{
		"api/v1": {config.apiHost + ":" + strconv.Itoa(config.apiPort)},
	})

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.PathPrefix("/api").Handler(proxy)

	n := negroni.Classic()
	n.UseHandler(r)

	log.Println("Starting web server on " + ":" + strconv.Itoa(config.defaultPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.defaultPort), n))
}
