package main

import (
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(res,
			`<table>
          <tr>
            <td style="text-align: center;">
              <img src="https://developerculture.com/wp-content/uploads/2016/01/golang-die-cut-3x3-1000x1000.jpg" width="500">
            </td>
            <td style="text-align: center;">
              The Gopher wants to:
              <a href="/party/">Party</a>
              <a href="/work/">Work</a>
            </td>
          </tr>
        </table>`)
	})

	mux.HandleFunc("/party/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(res, `<img src="http://gophergala.com/assets/img/fancy_gopher_renee.jpg" width="500">`)
	})

	mux.HandleFunc("/work/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(res, `<img src="http://d1ce5ur0bponn9.cloudfront.net/images/tech-blog/go-gopher-mascot-gotools.jpg" width="500">`)
	})

	http.ListenAndServe(":50000", mux)
}
