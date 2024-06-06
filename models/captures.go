package models

import (
	"time"

	"gorm.io/gorm"
)

type Capture struct {
	gorm.Model
	ServerID    uint      `json:"serverId" gorm:"column:server_id;type:uint;not null"`
	Server      Server    `json:"server" `
	GameID      uint      `json:"gameId" gorm:"column:game_id;type:uint;not null"`
	Game        Game      `json:"game" `
	RoundID     uint      `json:"roundId" gorm:"column:round_id;type:uint;not null"`
	Round       Round     `json:"round" `
	PlayerID    uint      `json:"playerId" gorm:"column:player_id;type:uint;not null"`
	Player      Player    `json:"player" `
	IsDestroy   uint      `json:"isDestroy" gorm:"column:is_destroy;type:uint;not null;default:0"`
	Objective   uint      `json:"objective" gorm:"column:objective;type:uint;not null"`
	CaptureTime time.Time `json:"captureTime" gorm:"column:capture_time;type:timestamp;not null"`
}
