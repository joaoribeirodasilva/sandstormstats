package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	Name      string    `gorm:"column:name;type:string;size:255;not null"`
	SteamID   string    `gorm:"column:steam_id;type:string;size:255;not null"`
	Kills     uint      `gorm:"column:kills;type:uint;not null;default:0"`
	Deaths    uint      `gorm:"column:deaths;type:uint;not null;default:0"`
	TeamKills uint      `gorm:"column:team_kills;type:uint;not null;default:0"`
	Games     uint      `gorm:"column:games;type:uint;not null;default:0"`
	Rounds    uint      `gorm:"column:rounds;type:uint;not null;default:0"`
	FirstSeen time.Time `gorm:"column:first_seen;type:timestamp"`
	LastSeen  time.Time `gorm:"column:last_seen;type:timestamp"`
}
