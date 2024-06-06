package log_parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/joaoribeirodasilva/sandstormstats/models"
	"gorm.io/gorm"
)

type LogRound struct {
	RoundID   *uint
	Round     uint
	WinTeam   *uint
	StartTime time.Time
	EndTime   *time.Time
}

func (p *LogParser) round(logLine string) error {

	isSR := p.regexps["is_start_round"].MatchString(logLine)
	if !isSR {
		isER := p.regexps["is_end_round"].MatchString(logLine)
		if !isER {
			return nil
		}
		if err := p.parseRound(logLine, false); err != nil {
			return err
		}
	} else {
		if err := p.parseRound(logLine, true); err != nil {
			return err
		}
	}

	return nil
}

func (p *LogParser) parseRound(str string, start bool) error {

	result := p.regexps["round_number"].FindStringSubmatch(str)
	if len(result) == 0 {
		return errors.New("no round number found")
	}

	temp := strings.TrimSpace(strings.Replace(result[0], "Round", "", 1))
	r, err := strconv.Atoi(temp)
	if err != nil {
		return err
	}

	var round *models.Round
	if start {

		if r == 1 {

			game := &models.Game{
				ServerID:  p.CurrentGame.Server.ID,
				MapID:     &p.CurrentGame.Map.ID,
				ModeID:    &p.CurrentGame.Mode.ID,
				StartTime: &p.CurrentTime,
			}

			if err := p.db.Conn.Create(game).Error; err != nil {
				return err
			}

			p.CurrentGame.Game = game
		}

		round = &models.Round{
			ServerID:  p.CurrentGame.Server.ID,
			GameID:    p.CurrentGame.Game.ID,
			Round:     uint(r),
			WinTeamID: nil,
			StartTime: p.CurrentTime,
			EndTime:   nil,
		}

		if err := p.db.Conn.Create(round).Error; err != nil {
			return err
		}

		p.CurrentGame.Round = round

	} else {
		round := &models.Round{}
		if err := p.db.Conn.Where(&models.Round{ServerID: p.CurrentGame.Server.ID, GameID: p.CurrentGame.Game.ID, Round: uint(r)}).First(round).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("end of round %d without a start, in game %d, sverver %s", r, p.CurrentGame.Game.ID, p.CurrentGame.Server.Key)
			}
			return errors.New("failed to query game round from database")
		}

		result = p.regexps["round_end_team"].FindStringSubmatch(str)
		if len(result) == 0 {
			return errors.New("no round team found")
		}
		temp := strings.TrimSpace(strings.Replace(result[0], "Team", "", 1))
		teamInt, err := strconv.Atoi(temp)
		team := uint(teamInt)
		if err != nil {
			return err
		}

		round.EndTime = &p.CurrentTime
		round.WinTeamID = &team

		if err := p.db.Conn.Save(round).Error; err != nil {
			return err
		}
		p.CurrentGame.Round = round

	}

	return nil
}
