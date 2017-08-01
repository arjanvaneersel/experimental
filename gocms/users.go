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

	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username       string `gorm:"not null; unique_index"`
	First_name     string
	Last_name      string
	Email          string `gorm:"not null; unique_index"`
	Password       string `gorm:"not null"`
	Pages          []Page
	ChangePassword bool
	IsAdmin        bool // System administrator
}

func (o *User) SetPassword(password string) error {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password+settings.Pepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	o.Password = string(hpass)
	return nil
}

type UserManager interface {
	Authenticate(q string, p string) (*User, error)
	Get(q string) (*User, error)
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint)
}

type UserGorm struct {
	*gorm.DB
}

func (o *UserGorm) Authenticate(q string, p string) (*User, error) {
	user := &User{}
	err := o.DB.Where("username = ? or email = ?", q, q).First(user).Error
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p+settings.Pepper))
	if err != nil {
		fmt.Println("Passoword error")
		return &User{}, err
	}

	return user, nil
}

func (o *UserGorm) Get(q string) (*User, error) {
	user := &User{}
	err := o.DB.Where(q).First(user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (o *UserGorm) GetByID(id uint) (*User, error) {
	user := &User{}
	err := o.DB.Where("id = ?", id).First(user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (o *UserGorm) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := o.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (o *UserGorm) Create(user *User) error {
	return o.DB.Create(user).Error
}

func (o *UserGorm) Update(user *User) error {
	return o.DB.Save(user).Error
}

func (o *UserGorm) Delete(id uint) {

}

func NewUserGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{db}
}

type Users struct {
	UserManager
}

func UserController(m UserManager) *Users {
	return &Users{
		UserManager: m,
	}
}
