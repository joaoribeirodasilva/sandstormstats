package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name;type:string;size:255;not null"`
}
