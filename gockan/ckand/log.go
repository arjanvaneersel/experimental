package main

import (
	"fmt"
	"log"
	"net/http"
)

// log the supplied http request
func LogRequest(req *http.Request) {
	var proto string
	if req.TLS == nil {
		proto = "http://"
	} else {
		proto = "https://"
	}
	message := fmt.Sprintf("%s (%s) - %s %s%s%s - %s",
		req.RemoteAddr, req.UserAgent, req.Method, proto, req.Host, req.URL.Path,
		req.Referer)
	log.Print(message)
}
