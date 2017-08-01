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
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

var users *Users
var sites *Sites
var pages *Pages

var loginView *View
var changePasswordView *View

func InitGorm(db *gorm.DB) {
	var err error
	var user = &User{}
	var page = &Page{}
	var site = &Site{}

	fmt.Println("Running migrations")
	//doGormMigrations(db)
	doGormReset(db)

	user, err = users.Get("username = 'admin'")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Creating a new administrator")

			user = &User{
				Username:       "admin",
				Email:          "info@balkan.tech",
				ChangePassword: true,
				IsAdmin:        true,
			}

			err = user.SetPassword("admin")
			if err != nil {
				panic(err)
			}

			err = users.Create(user)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	fmt.Println("Running site migrations")
	siteRoles := SiteRoleController(NewSiteRoleGorm(db))
	userSiteRoles := UserSiteRoleController(NewUserSiteUserRoleGorm(db))

	site, err = sites.GetByURL("localhost")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Creating a new site")

			site = &Site{
				Name: "Test site",
				URL:  "localhost",
			}
			err = sites.Create(site)
			if err != nil {
				panic(err)
			}

			fmt.Println("Creating the administrator role")
			adminrole := &SiteRole{
				Name:          "Administrator",
				Administrator: true,
				SiteID:        site.ID,
			}
			err = siteRoles.Create(adminrole)
			if err != nil {
				panic(err)
			}

			fmt.Println("Connecting administrator account to the default site")
			userrole := &UserSiteRole{
				SiteID: site.ID,
				UserID: user.ID,
				RoleID: adminrole.ID,
			}
			err = userSiteRoles.Create(userrole)
			if err != nil {
				panic(err)
			}

		} else {
			panic(err)
		}
	}

	page, err = pages.Get(1)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Creating a new root page")

			page = &Page{
				Slug:   "/",
				Title:  "Welcome to GoCMS!",
				Body:   "<h1>Welcome to GoCMS!</h1>",
				SiteID: site.ID,
			}

			err = pages.Create(page)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func SetupBaseViews() error {
	var err error

	loginView, err = NewView("base", "templates/base/login.gohtml")
	if err != nil {
		return err
	}

	changePasswordView, err = NewView("base", "templates/base/changepassword.gohtml")
	if err != nil {
		return err
	}

	return nil
}

func StartGormDB() (*gorm.DB, error) {
	if settings.DBType == "PG" {
		db, err := getGormPGConnection()
		if err != nil {
			return nil, err
		}
		db.LogMode(true)

		SetupBaseModelsViaGorm(db)
		return db, nil
	}
	return nil, errors.New("Invalid database type in settings. Valid options are: PG")
}

func SetupBaseModelsViaGorm(db *gorm.DB) {
	users = UserController(NewUserGorm(db))
	pages = PageController(NewPageGorm(db))
	sites = SiteController(NewSiteGorm(db))
}
