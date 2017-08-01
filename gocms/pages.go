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
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Page struct {
	gorm.Model
	Title      string
	Slug       string `gorm:"not null;unique_index"`
	Body       string
	ShowInMenu bool
	MenuTitle  string
	MenuWeight string
	SiteID     uint
}

type PageManager interface {
	Get(id uint) (*Page, error)
	GetBySlug(slug string, host string) (*Page, error)
	Create(page *Page) error
	Update(page *Page) error
	Delete(id uint)
}

type PageGorm struct {
	*gorm.DB
}

func (o *PageGorm) Get(id uint) (*Page, error) {
	page := &Page{}
	err := o.DB.First(&page, id).Error
	return page, err
}

func (o *PageGorm) GetBySlug(slug string, host string) (*Page, error) {
	var err error
	page := &Page{}
	site := &Site{}

	if settings.DefaultSite != 0 {
		err = o.DB.First(&site, settings.DefaultSite).Error
		if err != nil {
			return page, err
		}
	} else {
		err = o.DB.Where("URL = ?", host).First(&site).Error
		if err != nil {
			return page, err
		}
	}

	err = o.DB.Where("slug = ? and site_id = ?", slug, site.ID).First(page).Error
	return page, err
}

func (o *PageGorm) Create(page *Page) error {
	return o.DB.Create(page).Error
}

func (o *PageGorm) Update(page *Page) error {
	return o.DB.Update(page).Error
}

func (o *PageGorm) Delete(id uint) {

}

func NewPageGorm(db *gorm.DB) *PageGorm {
	return &PageGorm{db}
}

type Pages struct {
	PageManager
}

func PageController(m PageManager) *Pages {
	return &Pages{
		PageManager: m,
	}
}
