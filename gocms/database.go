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
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func getGormPGConnection() (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		settings.DBHost, settings.DBPort, settings.DBUser, settings.DBPassword, settings.DBName)
	db, err = gorm.Open("postgres", connection)
	if err != nil {
		return db, err
	}
	return db, nil
}

func doGormMigrations(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Site{}, &SiteRole{}, &UserSiteRole{}, &Page{})
}

func doGormReset(db *gorm.DB) {
	db.DropTableIfExists(&User{}, &Site{}, &SiteRole{}, &UserSiteRole{}, &Page{})
	doGormMigrations(db)
}

func DBType(s string) bool {
	return strings.ToLower(settings.DBType) == strings.ToLower(s)
}
