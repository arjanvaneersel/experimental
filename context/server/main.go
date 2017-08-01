package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/arjanvaneersel/experimental/context/clog"
)

func main() {
	http.HandleFunc("/", clog.Decorate(handler))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clog.Println(ctx, "Handler started")
	defer clog.Println(ctx,"Handler ended")

	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "Hello gopher!")
	case <-ctx.Done():
		clog.Println(ctx, ctx.Err().Error())
		http.Error(w, ctx.Err().Error(), http.StatusInternalServerError)
	}
}


