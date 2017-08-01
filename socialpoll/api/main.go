package main

import (
	"context"
	"net/http"
	"gopkg.in/mgo.v2"
	"flag"
	"log"
)

type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

func APIKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(contextKeyAPIKey).(string)
	return key, ok
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyAPIKey, key)
		fn(w, r.WithContext(ctx))
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

type Server struct {
	db *mgo.Session
}

func main(){
	var (
		addr = flag.String("addr", ":8080", "Endpoint address")
		mongo = flag.String("mongo", "localhost", "MongoDB address")
	)

	log.Println("Dialing mongo ", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("Failed to connect to mongodb: ", err)
	}
	defer db.Close()

	s := &Server{
		db: db,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withAPIKey(s.handlePolls)))
	log.Println("Starting web server on ", *addr)
	http.ListenAndServe(*addr, mux)
	log.Println("Stopping...")

}
