package models

import (
	"time"
)

type Exercise struct {
	ID                  uint `json:"id" gorm:"primarykey"`
	CategoryID          uint
	Category            Setting `gorm:"foreignKey:CategoryID"`
	SecondaryCategoryID *uint
	SecondaryCategory   Setting `gorm:"foreignKey:SecondaryCategoryID"`
	UserID              uint
	User                User
	Results []Result
	CreatedAt time.Time
}