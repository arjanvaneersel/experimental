package main

import (
  "log"
  "net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-type", "text/html")
  w.Write([]byte("<h1>Secure connection</h1>"))
}

func main() {
  http.HandleFunc("/", handler)

  err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
  if err != nil {
    log.Fatal(err)
  }
}
