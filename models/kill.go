package models

import (
	"time"

	"gorm.io/gorm"
)

type Kill struct {
	gorm.Model
	ServerID uint      `json:"serverId" gorm:"column:server_id;type:uint;not null"`
	Server   Server    `json:"server" `
	GameID   uint      `json:"gameId" gorm:"column:game_id;type:uint;not null"`
	Game     Game      `json:"game" `
	RoundID  uint      `json:"roundId" gorm:"column:round_id;type:uint;not null"`
	Round    Round     `json:"round" `
	MapID    uint      `json:"mapId" gorm:"column:map_id;type:uint;not null"`
	Map      Map       `json:"map" `
	PlayerID uint      `json:"playerId" gorm:"column:player_id;type:uint;not null"`
	Player   Player    `json:"player" `
	RoleID   uint      `json:"roleId" gorm:"column:role_id;type:uint;not null"`
	Role     Role      `json:"role" `
	Weapon   string    `json:"weapon" gorm:"column:weapon;type:string;size:255;not null"`
	KillTime time.Time `json:"killTime" gorm:"column:kill_time;type:timestamp;not null"`
	IsShared uint      `json:"isShared" gorm:"column:is_shared;type:uint;not null;default:0"`
	Killed   uint      `json:"killed" gorm:"column:killed;type:uint;not null;default:0"`
}
