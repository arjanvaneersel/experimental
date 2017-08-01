/* GoCMS: A basic Golang CMS system
Copyright (C) 2016, "Balkan C & T" OOD

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"flag"
)

func NewRouter() *mux.Router {
	var dir string

	flag.StringVar(&dir, "dir", "./static", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc(settings.Login_URL, loginPage).Methods("GET")
	r.HandleFunc(settings.Login_URL, loginHandler).Methods("POST")
	r.HandleFunc(settings.Password_URL, RequireLogin(changePasswordPage)).Methods("GET")
	r.HandleFunc(settings.Password_URL, RequireLogin(changePasswordHandler)).Methods("POST")
	r.HandleFunc(settings.Logout_URL, RequireLogin(logoutHandler))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		showPage(w, r)
	})
	r.HandleFunc("/{language:[a-z]{2}}", showPage)
	r.HandleFunc("/{language:[a-z]{2}}/{slug:[a-z]+}", showPage)
	r.HandleFunc("/{slug:[a-z]+}", showPage)
	r.NotFoundHandler = http.HandlerFunc(notfound)

	return r
}
