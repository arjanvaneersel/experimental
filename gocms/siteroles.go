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

type SiteRole struct {
	gorm.Model
	SiteID        uint
	Name          string `gorm:"unique_index"`
	Administrator bool
	MayCreate     bool
	MayEdit       bool
	MayModerate   bool
	MayComment    bool
}

type SiteRoleManager interface {
	Get(id uint) *SiteRole
	GetByName(name string) (*SiteRole, error)
	Create(site *SiteRole) error
	Update(site *SiteRole) error
	Delete(id uint)
}

type SiteRoleGorm struct {
	*gorm.DB
}

func (o *SiteRoleGorm) Get(id uint) *SiteRole {
	return nil
}

func (o *SiteRoleGorm) GetByName(name string) (*SiteRole, error) {
	role := &SiteRole{}
	err := o.DB.Where("name = ?", name).First(role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (o *SiteRoleGorm) Create(role *SiteRole) error {
	return o.DB.Create(role).Error
}

func (o *SiteRoleGorm) Update(role *SiteRole) error {
	return o.DB.Update(role).Error
}

func (o *SiteRoleGorm) Delete(id uint) {

}

func NewSiteRoleGorm(db *gorm.DB) *SiteRoleGorm {
	return &SiteRoleGorm{db}
}

type SiteRoles struct {
	SiteRoleManager
}

func SiteRoleController(m SiteRoleManager) *SiteRoles {
	return &SiteRoles{
		SiteRoleManager: m,
	}
}
