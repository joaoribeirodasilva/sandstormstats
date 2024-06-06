package models

import (
	"time"

	"gorm.io/gorm"
)

type ServerLog struct {
	gorm.Model
	ServerID  uint       `json:"serverId" gorm:"column:server_id;type:uint;not null"`
	Server    Server     `json:"server"`
	File      string     `json:"file" gorm:"column:file;type:string;size:255;not null"`
	Success   uint       `json:"success" gorm:"column:success;type:uint;size:1;not null:dedfault:0"`
	Message   string     `json:"message" gorm:"column:message;type:string;size:65535"`
	StartTime time.Time  `json:"startTime" gorm:"column:start_time;type:timestamp;not null"`
	EndTime   *time.Time `json:"sendTime" gorm:"column:end_time;type:timestamp;"`
}
