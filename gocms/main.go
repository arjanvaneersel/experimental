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

import "log"

//var db *gorm.DB

var settings = Settings{
	Debug:        true,
	ServeHost:    "",
	ServePort:    8000,
	DBType:       "PG",
	DBHost:       "localhost",
	DBPort:       5432,
	DBUser:       "gocms",
	DBPassword:   "gocms",
	DBName:       "gocms",
	DefaultSite:  1,
	Pepper:       "t44etetETRTY%$t09iyt0)yi0945yU84584Y84595909iRg5p54)",
	Csrf:         "re4544t5y6htrrt5yu7ik86kjsherge5tyu76e56tergre",
	Login_URL:    "/login",
	Logout_URL:   "/logout",
	Password_URL: "/password",
}

func main() {
	if DBType("PG") {
		db, err := StartGormDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		InitGorm(db)
	} else {
		log.Fatal("Invalid Database type in settings. Please check settings.DBType.")
	}

	err := StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
