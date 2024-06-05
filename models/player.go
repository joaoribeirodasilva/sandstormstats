package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	Name      string     `json:"name" gorm:"column:name;type:string;size:255;not null"`
	SteamID   string     `json:"steamId" gorm:"column:steam_id;type:string;size:255;not null"`
	Kills     uint       `json:"kills" gorm:"column:kills;type:uint;not null;default:0"`
	Deaths    uint       `json:"deaths" gorm:"column:deaths;type:uint;not null;default:0"`
	TeamKills uint       `json:"teamKills" gorm:"column:team_kills;type:uint;not null;default:0"`
	Games     uint       `json:"games" gorm:"column:games;type:uint;not null;default:0"`
	Rounds    uint       `json:"rounds" gorm:"column:rounds;type:uint;not null;default:0"`
	FirstSeen time.Time  `json:"firstSeen" gorm:"column:first_seen;type:timestamp"`
	LastSeen  *time.Time `json:"lastSeen" gorm:"column:last_seen;type:timestamp"`
}
