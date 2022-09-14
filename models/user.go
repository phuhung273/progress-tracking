package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique;size:50"`
	Password string `gorm:"not null;size:250"`
}