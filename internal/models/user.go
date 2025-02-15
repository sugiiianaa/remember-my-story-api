package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string
	FullName string
	Password string
}
