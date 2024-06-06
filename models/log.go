package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	ServerID uint   `json:"serverId" gorm:"column:server_id;type:uint;not null"`
	Server   Server `json:"server" `
	Entry    string `json:"entry" gorm:"column:entry;type:string;size:65535"`
}
