package models

import (
	"time"

	"gorm.io/gorm"
)

type Round struct {
	gorm.Model
	ServerID  uint       `json:"serverId" gorm:"column:server_id;type:uint;not null"`
	Server    Server     `json:"server"`
	GameID    uint       `json:"gameId" gorm:"column:game_id;type:uint;not null"`
	Game      Game       `json:"game"`
	WinTeamID *uint      `json:"winTeamId" gorm:"column:win_team_id;type:uint;not null"`
	WinTeam   Team       `json:"winTeam"`
	StartTime time.Time  `json:"startTime" gorm:"column:start_time;type:timestamp;not null"`
	EndTime   *time.Time `json:"sendTime" gorm:"column:end_time;type:timestamp;"`
}
