package models

import "gorm.io/gorm"

type Map struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name;type:string;size:255;not null"`
	File string `json:"file" gorm:"column:file;type:string;size:255;not null"`
}
