package log_parser

import (
	"errors"
	"strings"

	"github.com/joaoribeirodasilva/sandstormstats/models"
	"gorm.io/gorm"
)

type Kill struct {
	Players []Player
	Dead    Player
	Weapon  string
}

func (p *LogParser) kill(logLine string) error {

	is, err := p.isKill(logLine)
	if err != nil {
		return err
	}
	if !is {
		return nil
	}

	return nil
}

func (p *LogParser) isKill(str string) (bool, error) {

	isK := p.regexps["is_kill"].MatchString(str)
	if !isK {
		return false, nil
	}

	entry := p.regexps["has_date"].ReplaceAllString(str, "")

	idx := strings.Index(entry, "Display:")
	idx = idx + len("Display:")
	entry = entry[idx:]

	idxWith := strings.Index(entry, "with")
	idx = idxWith + len("with")
	weapon := entry[idx:]
	weapon = strings.TrimSpace(weapon)

	idxKilled := strings.Index(entry, "killed")
	idx = idxKilled + len("killed")
	killed := entry[idx:idxWith]

	playersString := entry[:idxKilled]
	players, err := p.parsePlayers(playersString, false)
	if err != nil {
		return true, err
	}
	dead, err := p.parsePlayer(killed)
	if err != nil {
		return true, err
	}

	for _, item := range *players {

		if item.SteamID == "" && dead.SteamID == "" {
			continue
		}

		isDeath := 0
		isKill := 0
		isTk := 0
		playerId := item.SteamID
		playerName := item.Name
		killerName := item.Name
		deadName := dead.Name
		if item.SteamID == "" && dead.SteamID != "" { //death
			isDeath = 1
			isKill = 0
			isTk = 0
			playerId = dead.SteamID
			playerName = dead.SteamID
		} else if item.SteamID != "" && item.SteamID == dead.SteamID { //suicide
			isDeath = 1
			isKill = 1
			isTk = 0
		} else if dead.SteamID != "" && item.SteamID != "" && item.SteamID != dead.SteamID { //team kill
			isDeath = 0
			isKill = 1
			isTk = 1
		}

		player := &models.Player{}
		if err := p.db.Conn.Where(&models.Player{SteamID: playerId}).First(player).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				player.Name = playerName
				player.SteamID = playerId
				player.FirstSeen = p.CurrentTime
				player.LastSeen = &p.CurrentTime
				if err := p.db.Conn.Create(player).Error; err != nil {
					return true, err
				}
			} else {
				return true, err
			}
		} else {
			player.LastSeen = &p.CurrentTime
			if err := p.db.Conn.Save(player).Error; err != nil {
				return true, err
			}
		}

		gamePlayer := &models.GamePlayer{}
		if err := p.db.Conn.Where(&models.GamePlayer{ServerID: p.CurrentGame.Server.ID, GameID: p.CurrentGame.Game.ID, PlayerID: player.ID}).First(gamePlayer).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				gamePlayer.ServerID = p.CurrentGame.Server.ID
				gamePlayer.GameID = p.CurrentGame.Game.ID
				gamePlayer.PlayerID = player.ID
				gamePlayer.StartTime = p.CurrentTime
				gamePlayer.TeamID = item.Team
				if err := p.db.Conn.Create(gamePlayer).Error; err != nil {
					return true, err
				}
			} else {
				return true, err
			}
		}

		var mapId *uint = nil
		if p.CurrentGame.Map != nil {
			mapId = &p.CurrentGame.Map.ID
		}

		var shared uint = 0
		if len(*players) > 1 {
			shared = 1
		}

		var suicide uint = 0
		if isDeath == 1 && isKill == 1 {
			suicide = 1
		}

		kill := &models.Kill{
			ServerID: p.CurrentGame.Server.ID,
			GameID:   p.CurrentGame.Game.ID,
			RoundID:  p.CurrentGame.Round.ID,
			MapID:    mapId,
			PlayerID: player.ID,
			KillTime: p.CurrentTime,
			Weapon:   weapon,
			IsShared: shared,
			Killed:   uint(isDeath),
			Suicide:  suicide,
			TeamKill: uint(isTk),
			Killer:   killerName,
			Dead:     deadName,
		}

		if err := p.db.Conn.Create(kill).Error; err != nil {
			return true, err
		}
	}

	return true, nil
}
