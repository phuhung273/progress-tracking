package models

type Setting struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;size:50"`
	Type string `gorm:"not null;size:50"`
}