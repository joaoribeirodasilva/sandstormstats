package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	ServerID  uint       `json:"serverId" gorm:"column:server_id;type:uint;not null"`
	Server    Server     `json:"server" `
	MapID     *uint      `json:"mapId" gorm:"column:map_id;type:uint"`
	Map       *Map       `json:"map" `
	ModeID    *uint      `json:"modeId" gorm:"column:mode_id;type:uint"`
	Mode      *Mode      `json:"mode" `
	StartTime *time.Time `json:"startTime" gorm:"column:start_time;type:timestamp;not null"`
	EndTime   *time.Time `json:"endType" gorm:"column:end_time;type:timestamp"`
}
