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
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
)

func StartServer() error {
	err := SetupBaseViews()
	if err != nil {
		log.Fatal(err)
	}

	r := NewRouter()
	CSRF := csrf.Protect([]byte(settings.Csrf), csrf.Secure(false))
	server := settings.ServeHost + ":" + strconv.Itoa(int(settings.ServePort))
	log.Printf("Starting the server on %s", server)
	err = http.ListenAndServe(server, CSRF(r))
	if err != nil {
		return err
	}
	return nil
}
