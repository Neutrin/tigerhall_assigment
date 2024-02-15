package model

import "gorm.io/gorm"

type User struct {
	UserName string `json:"username"`
	//making it unique for login
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	gorm.Model
	//role and other fields can be added in future
}
