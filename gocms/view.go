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

	"github.com/gorilla/csrf"
)

func NewView(layout string, files ...string) (*View, error) {
	files = append(files,
		"templates/includes/navbar.gohtml",
		"templates/layouts/"+layout+".gohtml",
	)
	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return &View{
		Template: t,
		Layout:   layout,
	}, nil
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) Execute(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	activeUser, err := getUser(r)
	if err != nil {
		fmt.Println(err)
	}

	d := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"ActiveUser":     activeUser,
		"Alerts":         alerts,
	}

	if data != nil {
		for k, v := range data {
			d[k] = v
		}
	}

	clearAlerts()

	w.Header().Set("Content-Type", "text/html")
	v.Template.ExecuteTemplate(w, v.Layout, d)
}
