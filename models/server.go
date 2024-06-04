package models

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	Key  string `gorm:"column:key;type:string;size:255;not null"`
	Name string `gorm:"column:name;type:string;size:255;not null"`
}
