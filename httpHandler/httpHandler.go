package httpHandler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kr/pretty"
	"github.com/urfave/negroni"
	"github.com/xu354cjo1008/eatingFinder/geography/geocoding"
)

func HttpGet(request string) ([]byte, error) {
	resp, err := http.Get(request)
	if err != nil {
		return nil, errors.New("http.get failed")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read http response failed")
	}

	return body, nil
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Home")
}

func apiHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Api Handler")
}

func apiGeocodeHandler(rw http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(rw, "Api Geocode Handler")

	vars := mux.Vars(r)

	geocode := geocoding.NewGeocode("AIzaSyDJXVVPUtvmRDcBN4nTPNVAI26cUzOaztw", "en")

	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lng, _ := strconv.ParseFloat(vars["lng"], 64)

	city, err := geocode.GetCityByLatlng(lat, lng)

	if err != nil {
		log.Println("error: ", err)
	}

	fmt.Fprintln(rw, city)
}

func middleware1(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	pretty.Println("middleware handler")
	next(rw, r)
}

func RunServer(port int) {

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(false)
	apiRouter.HandleFunc("/", apiHandler)
	apiRouter.HandleFunc("/getCity/{lat}/{lng}", apiGeocodeHandler)

	r.PathPrefix("/api").Handler(apiRouter)

	n := negroni.Classic()
	n.UseHandler(r)

	log.Println("Starting server on " + ":" + strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), n))
}
