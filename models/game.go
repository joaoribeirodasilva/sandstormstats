package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	ServerID  uint      `gorm:"column:server_id;type:uint;not null"`
	StartTime time.Time `gorm:"column:start_time;type:timestamp;not null"`
	EndTime   time.Time `gorm:"column:end_time;type:timestamp"`
	MapID     uint      `gorm:"column:map_id;type:uint;not null"`
	ModeID    uint      `gorm:"column:mode_id;type:uint;not null"`
}
