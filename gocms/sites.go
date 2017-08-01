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

type Site struct {
	gorm.Model
	Name string
	URL  string `gorm:"not null; unique_index"`
}

type SiteManager interface {
	Get(id uint) *Site
	GetByURL(slug string) (*Site, error)
	Create(site *Site) error
	Update(site *Site) error
	Delete(id uint)
}

type SiteGorm struct {
	*gorm.DB
}

func (o *SiteGorm) Get(id uint) *Site {
	return nil
}

func (o *SiteGorm) GetByURL(URL string) (*Site, error) {
	site := &Site{}
	err := o.DB.Where("URL = ?", URL).First(site).Error
	if err != nil {
		return site, err
	}
	return site, nil
}

func (o *SiteGorm) Create(site *Site) error {
	return o.DB.Create(site).Error
}

func (o *SiteGorm) Update(site *Site) error {
	return o.DB.Update(site).Error
}

func (o *SiteGorm) Delete(id uint) {

}

func NewSiteGorm(db *gorm.DB) *SiteGorm {
	return &SiteGorm{db}
}

type Sites struct {
	SiteManager
}

func SiteController(m SiteManager) *Sites {
	return &Sites{
		SiteManager: m,
	}
}
