package models

import "gorm.io/gorm"

type Mode struct {
	gorm.Model
	Name  string `json:"name" gorm:"column:name;type:string;size:255;not null"`
	Class string `json:"class" gorm:"column:class;type:string;size:255;not null"`
}
