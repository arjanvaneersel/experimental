package main

import (
	"net/http"
	"encoding/json"
	"github.com/arjanvaneersel/meander"
	"strings"
	"strconv"
)

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w,r)
	}
}

func main() {
	meander.APIKey = "AIzaSyACWouN1S9Yzw9IkLqbRj8FuI5OlbzpNZI"
	http.HandleFunc("/journeys", cors(func(w http.ResponseWriter, r *http.Request){
		respond(w, r, meander.Journeys)
	}))

	http.HandleFunc("/recommendations", cors(func(w http.ResponseWriter, r *http.Request) {
		q := &meander.Query{Journey: strings.Split(r.URL.Query().Get("journey"), "|")}
		var err error

		q.Lat, err = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		q.Lng, err = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		q.Radius, err = strconv.Atoi(r.URL.Query().Get("radius"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		q.CostRangeStr =  r.URL.Query().Get("cost")
		places := q.Run()
		respond(w, r, places)
	}))


	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}
	return json.NewEncoder(w).Encode(publicData)
}
