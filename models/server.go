package models

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	Key         string  `json:"key" gorm:"column:key;type:string;size:255;not null"`
	Name        string  `json:"name" gorm:"column:name;type:string;size:255;not null"`
	LogDir      string  `json:"logDir" gorm:"column:log_dir;type:string;size:255"`
	LastLog     *string `json:"lastLog" gorm:"column:last_log;type:string;size:255"`
	LastLine    uint    `json:"lastLine" gorm:"column:last_line;type:uint;not null;default:0"`
	LastLogSize uint64  `json:"lastLogSize" gorm:"column:last_log_size;type:uint;not null;default:0"`
}
