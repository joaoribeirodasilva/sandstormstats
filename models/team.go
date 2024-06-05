package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name;type:string;size:50;not null"`
}
