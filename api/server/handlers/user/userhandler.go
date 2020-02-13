package user

import (
	"errors"
	"strings"
)

type User struct {
	Id           string   `json:"-" bson:"_id_"`
	Name         string   `json:"name" bson:"name"`
	Surname      string   `json:"surname" bson:"surname"`
	Workdays     int      `json:"workdays" bson:"workdays"`
	Weekdays     int      `json:"weekdays" bson:"weekdays"`
	FirstWorkDay string   `json:"firstwork" bson:"firstwork"`
	UserCreated  int64    `json:"-" bson:"created"`
	Email        string   `json:"email" bson:"email"`
	Phone        string   `json:"phone" bson:"phone"`
	Groups       []string `json:"groups" bson:"groups"`
}

type UserAuth struct {
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
}

func (ua UserAuth) HasRequiredFields() error {
	if len(ua.Name) == 0 {
		return errors.New("Name is emprty")
	} else if len(ua.Email) == 0 || !strings.Contains(ua.Email, "@") {
		return errors.New("Email not valid")
	} else if len(ua.Phone) < 11 {
		return errors.New("Phone is invalid")
	} else if len(ua.Password) < 5 {
		return errors.New("Password is invalid or missing")
	} else {
		return nil
	}
}

func (u User) HasRequiredFields() error {

	if len(u.Name) == 0 {
		return errors.New("Name is empty")
	} else if len(u.Surname) == 0 {
		return errors.New("Surname is empty")
	} else if u.Workdays < 1 || u.Weekdays < 1 {
		return errors.New("Graph not filled correctly")
	} else if u.Email == "" {
		return errors.New("Email is empty")
	} else if u.Phone == "" {
		return errors.New("Phone is empty")
	} else {
		return nil
	}
}
