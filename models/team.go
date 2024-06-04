package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name string `gorm:"column:name;type:string;size:50;not null"`
}
