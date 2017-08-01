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

type UserSiteRole struct {
	gorm.Model
	SiteID uint
	UserID uint
	RoleID uint
}

type UserSiteRoleManager interface {
	Get(id uint) (*UserSiteRole, error)
	Create(role *UserSiteRole) error
	Update(role *UserSiteRole) error
	Delete(id uint)
}

type UserSiteGorm struct {
	*gorm.DB
}

func (o *UserSiteGorm) Get(id uint) (*UserSiteRole, error) {
	role := &UserSiteRole{}
	err := o.DB.First(role, id).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (o *UserSiteGorm) Create(role *UserSiteRole) error {
	return o.DB.Create(role).Error
}

func (o *UserSiteGorm) Update(role *UserSiteRole) error {
	return o.DB.Update(role).Error
}

func (o *UserSiteGorm) Delete(id uint) {

}

func NewUserSiteUserRoleGorm(db *gorm.DB) *UserSiteGorm {
	return &UserSiteGorm{db}
}

type UserSiteRoles struct {
	UserSiteRoleManager
}

func UserSiteRoleController(m UserSiteRoleManager) *UserSiteRoles {
	return &UserSiteRoles{
		UserSiteRoleManager: m,
	}
}
