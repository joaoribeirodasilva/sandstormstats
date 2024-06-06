package models

import (
	"time"

	"gorm.io/gorm"
)

type GamePlayer struct {
	gorm.Model
	ServerID  uint       `json:"serverId" gorm:"column:server_id;type:uint;not null"`
	Server    Server     `json:"server" `
	GameID    uint       `json:"gameId" gorm:"column:game_id;type:uint;not null"`
	Game      Game       `json:"game" `
	TeamID    uint       `json:"teamId" gorm:"column:team_id;type:uint;not null"`
	Team      Team       `json:"team" `
	PlayerID  uint       `json:"playerId" gorm:"column:player_id;type:uint;not null"`
	Player    Player     `json:"player" `
	StartTime time.Time  `json:"startTime" gorm:"column:start_time;type:timestamp;not null"`
	EndTime   *time.Time `json:"endTime" gorm:"column:end_time;type:timestamp;not null"`
}
