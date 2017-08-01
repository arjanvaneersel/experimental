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

	"log"

	"github.com/gorilla/securecookie"
)

type ActiveUser struct {
	ID       uint
	Username string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func createSession(u *User, w http.ResponseWriter) {
	v := map[string]interface{}{
		"id":       u.ID,
		"username": u.Username,
	}

	if encoded, err := cookieHandler.Encode("session", v); err == nil {
		c := &http.Cookie{
			Name:  "gocms",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, c)
	}
}

func destroySession(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:   "gocms",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
}

func getUser(r *http.Request) (user ActiveUser, err error) {
	c, err := r.Cookie("gocms")

	if err != nil {
		return user, err
	}

	cValue := make(map[string]interface{})
	if err = cookieHandler.Decode("session", c.Value, &cValue); err != nil {
		return user, err
	}
	return ActiveUser{ID: cValue["id"].(uint), Username: cValue["username"].(string)}, nil
}

var RequireLoginRedirectTo string = "/login"

func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUser(r)
		if err != nil {
			alerts = append(alerts, Alert{Title: "Warning", Class: "warning", Message: "You need to login to access that page."})
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
		}

		_, err = users.GetByID(user.ID)
		if err != nil {
			alerts = append(alerts, Alert{Title: "Warning", Class: "warning", Message: "You need to login to access that page."})
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
		}
		next(w, r)
	}
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUser(r)
		if err != nil {
			alerts = append(alerts, Alert{Title: "Warning", Class: "warning", Message: "You need to login to access that page."})
			http.Redirect(w, r, RequireLoginRedirectTo, http.StatusTemporaryRedirect)
			return
		}

		u, err := users.GetByID(user.ID)
		if !u.IsAdmin || err != nil {
			log.Printf("%s: Non-admin access attempt", r.URL.Path)
			http.Error(w, "Access denied", http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}
