package models

import (
	"time"

	"gorm.io/datatypes"
)

type Exercise struct {
	ID                  uint `gorm:"primarykey"`
	CategoryID          uint
	Category            Setting `gorm:"foreignKey:CategoryID"`
	SecondaryCategoryID *uint
	SecondaryCategory   Setting `gorm:"foreignKey:SecondaryCategoryID"`
	UserID              uint
	User                User
	CreatedAt time.Time
	Result datatypes.JSON
}