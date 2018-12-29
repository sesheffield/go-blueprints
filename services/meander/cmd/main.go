package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/sesheffield/go-blueprints/services/meander"
)

var config struct {
	GMapsAPIKey string `required:"true" envconfig:"GMAPS_API_KEY"`
}

func init() {
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("unable to configure from environment: %v", err)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "The address of the application")
	flag.Parse()
	// Use all of those CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())
	meander.APIKey = config.GMapsAPIKey
	http.HandleFunc("/journeys", cors(func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	}))
	http.HandleFunc("/recommendations", cors(func(w http.ResponseWriter, r *http.Request) {
		q := &meander.Query{
			Journey: strings.Split(r.URL.Query().Get("journey"), "|"),
		}
		q.Lat, _ = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		q.Lng, _ = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
		q.Radius, _ = strconv.Atoi(r.URL.Query().Get("radius"))
		q.CostRangeStr = r.URL.Query().Get("cost")
		places := q.Run()
		respond(w, r, places)
	}))
	log.Printf("Listening on port %s...", *addr)
	http.ListenAndServe(*addr, http.DefaultServeMux)
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(&publicData)
}
