package log_parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/joaoribeirodasilva/sandstormstats/models"
	"gorm.io/gorm"
)

type Objective struct {
	Players   []Player
	Objective uint
	Team      uint
	Destroy   bool
}

func (p *LogParser) objective(logLine string) error {

	logLine = p.regexps["has_date"].ReplaceAllString(logLine, "")

	is, err := p.isCapture(logLine)
	if err != nil {
		return err
	}

	if !is {
		if err = p.isDestroy(logLine); err != nil {
			return err
		}
	}

	return nil
}

func (p *LogParser) isCapture(str string) (bool, error) {

	isC := p.regexps["is_capture"].MatchString(str)
	if !isC {
		return false, nil
	}

	if err := p.parseObjective(str, false); err != nil {
		return true, err
	}

	return true, nil
}

func (p *LogParser) isDestroy(str string) error {

	isD := p.regexps["is_destroy"].MatchString(str)
	if !isD {
		return nil
	}

	if err := p.parseObjective(str, true); err != nil {
		return err
	}

	return nil
}

func (p *LogParser) parseObjective(str string, destroy bool) error {

	result := p.regexps["objective_get_num"].FindStringSubmatch(str)
	if len(result) == 0 {
		return errors.New("no objective number found")
	}

	temp := strings.TrimSpace(strings.Replace(result[0], "Objective", "", 1))
	objNumber, err := strconv.Atoi(temp)
	if err != nil {
		return err
	}

	result = p.regexps["objective_get_team"].FindStringSubmatch(str)
	if len(result) == 0 {
		return errors.New("no objective team found")
	}

	temp = strings.TrimSpace(strings.Replace(result[0], "for team", "", 1))
	objTeam, err := strconv.Atoi(temp)
	if err != nil {
		return err
	}

	result = p.regexps["objective_get_players"].FindStringSubmatch(str)
	if len(result) == 0 {
		return errors.New("no objective players found")
	}
	temp = strings.TrimSpace(strings.Replace(result[0], "by", "", 1))
	temp = strings.ReplaceAll(temp, ",", " +")
	if temp == "" {
		return errors.New("no objective players found")
	}

	players, err := p.parsePlayers(temp, true)
	if err != nil {
		return err
	}

	for idx := range *players {
		if (*players)[idx].SteamID == "INVALID" {
			continue
		}
		player := &models.Player{}
		if err := p.db.Conn.Where(&models.Player{SteamID: (*players)[idx].SteamID}).First(player).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				player.Name = (*players)[idx].Name
				player.SteamID = (*players)[idx].SteamID
				player.FirstSeen = p.CurrentTime
				player.LastSeen = &p.CurrentTime
				if err := p.db.Conn.Create(player).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			player.LastSeen = &p.CurrentTime
			if err := p.db.Conn.Save(player).Error; err != nil {
				return err
			}
		}

		gamePlayer := &models.GamePlayer{}
		if err := p.db.Conn.Where(&models.GamePlayer{ServerID: p.CurrentGame.Server.ID, GameID: p.CurrentGame.Game.ID, PlayerID: player.ID}).First(gamePlayer).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				gamePlayer.ServerID = p.CurrentGame.Server.ID
				gamePlayer.GameID = p.CurrentGame.Game.ID
				gamePlayer.PlayerID = player.ID
				gamePlayer.StartTime = p.CurrentTime
				gamePlayer.TeamID = uint(objTeam)
				if err := p.db.Conn.Create(gamePlayer).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		capture := &models.Capture{}
		if err := p.db.Conn.Where(&models.Capture{ServerID: p.CurrentGame.Server.ID, GameID: p.CurrentGame.Game.ID, RoundID: p.CurrentGame.Round.ID, PlayerID: player.ID}).First(capture).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				capture.ServerID = p.CurrentGame.Server.ID
				capture.GameID = p.CurrentGame.Game.ID
				capture.RoundID = p.CurrentGame.Round.ID
				capture.PlayerID = player.ID
				capture.IsDestroy = 0
				capture.Objective = uint(objNumber)
				if destroy {
					capture.IsDestroy = 1
				}
				capture.CaptureTime = p.CurrentTime
				if err := p.db.Conn.Create(capture).Error; err != nil {
					return err
				}
			}
		}
		(*players)[idx].Team = uint(objTeam)
	}

	return nil
}
