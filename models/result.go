package models

type Result struct {
	ID         uint `gorm:"primarykey"`
	Value      uint
	CriteriaID uint
	Criteria   Setting `gorm:"foreignKey:CriteriaID"`
	ExerciseID uint
	Exercise   Exercise
}