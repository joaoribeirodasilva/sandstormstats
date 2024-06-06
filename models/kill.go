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
	MapID    *uint     `json:"mapId" gorm:"column:map_id;type:uint"`
	Map      Map       `json:"map" `
	PlayerID uint      `json:"playerId" gorm:"column:player_id;type:uint;not null"`
	Player   Player    `json:"player" `
	Weapon   string    `json:"weapon" gorm:"column:weapon;type:string;size:255;not null"`
	KillTime time.Time `json:"killTime" gorm:"column:kill_time;type:timestamp;not null"`
	IsShared uint      `json:"isShared" gorm:"column:is_shared;type:uint;not null;default:0"`
	Killed   uint      `json:"killed" gorm:"column:killed;type:uint;not null;default:0"`
	Suicide  uint      `json:"suicide" gorm:"column:suicide;type:uint;not null;default:0"`
	TeamKill uint      `json:"teamKill" gorm:"column:team_kill;type:uint;not null;default:0"`
	Killer   string    `json:"killer" gorm:"column:killer;type:string;size:255"`
	Dead     string    `json:"dead" gorm:"column:dead;type:string;size:255"`
}
