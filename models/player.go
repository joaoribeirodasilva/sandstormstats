package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	Name      string     `json:"name" gorm:"column:name;type:string;size:255;not null"`
	SteamID   string     `json:"steamId" gorm:"column:steam_id;type:string;size:255;not null"`
	FirstSeen time.Time  `json:"firstSeen" gorm:"column:first_seen;type:timestamp"`
	LastSeen  *time.Time `json:"lastSeen" gorm:"column:last_seen;type:timestamp"`
}
