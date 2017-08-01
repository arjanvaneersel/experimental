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
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h2>Oops... That page couldn't be found.</h2>")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	loginView.Execute(w, r, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	u, err := users.Authenticate(username, password)
	if err != nil {
		alerts = append(alerts, Alert{Title: "Error", Class: "danger", Message: "Invalid login"})
		http.Redirect(w, r, settings.Login_URL, 302)
		return
	}
	createSession(u, w)
	if u.ChangePassword {
		alerts = append(alerts, Alert{Title: "Warning", Class: "warning", Message: "You need to change your password."})
		http.Redirect(w, r, settings.Password_URL, 302)
		return
	}
	alerts = append(alerts, Alert{Title: "Success", Class: "success", Message: "You have succesfully logged in."})
	http.Redirect(w, r, "/", 302)
	return
}

func changePasswordPage(w http.ResponseWriter, r *http.Request) {
	changePasswordView.Execute(w, r, map[string]interface{}{"Url": settings.Password_URL})
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		a, _ := getUser(r)
		oldpassword := r.FormValue("oldpassword")
		password := r.FormValue("password")
		passwordc := r.FormValue("passwordc")

		u, err := users.Authenticate(a.Username, oldpassword)
		if err != nil {
			alerts = append(alerts, Alert{Title: "Error", Class: "danger", Message: "Old password is incorrect."})
			http.Redirect(w, r, settings.Password_URL, 302)
			return
		}

		if password != passwordc {
			alerts = append(alerts, Alert{Title: "Error", Class: "danger", Message: "Passwords are not the same."})
			http.Redirect(w, r, settings.Password_URL, 302)
			return
		}

		u.SetPassword(password)
		u.ChangePassword = false

		err = users.Update(u)
		if err != nil {
			alerts = append(alerts, Alert{Title: "Error", Class: "danger", Message: "Couldn't update the user."})
			http.Redirect(w, r, settings.Password_URL, 302)
			return
		}

		destroySession(w)
		alerts = append(alerts, Alert{Title: "Success", Class: "success", Message: "You have succesfully changed your password. Please login with your new password."})
		http.Redirect(w, r, "/", 302)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	destroySession(w)
	alerts = append(alerts, Alert{Title: "Success", Class: "success", Message: "You have succesfully logged out."})
	http.Redirect(w, r, "/", 302)
	return
}

func showPage(w http.ResponseWriter, r *http.Request) {
	var slug string

	host := strings.ToLower(strings.Split(r.Host, ":")[0])

	vars := mux.Vars(r)
	if vars["slug"] == "" {
		slug = "/"
	} else {
		slug = vars["slug"]
	}

	lang := vars["language"]
	fmt.Println(lang, slug)
	page, err := pages.GetBySlug(slug, host)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			notfound(w, r)
			return
		}
	}

	pageView, err := NewView("base", "templates/page.gohtml")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html")
	pageView.Execute(w, r, map[string]interface{}{
		"Title": page.Title,
		"Body":  template.HTML(page.Body),
	})
}
