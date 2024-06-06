package log_parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joaoribeirodasilva/sandstormstats/models"
	"gorm.io/gorm"
)

func (p *LogParser) mode(logLine string) error {

	isS := p.regexps["is_mode"].MatchString(logLine)
	if !isS {
		return nil
	}

	if err := p.parseMode(logLine); err != nil {
		return err
	}

	return nil
}

func (p *LogParser) parseMode(str string) error {

	result := p.regexps["mode_class"].FindStringSubmatch(str)
	if len(result) == 0 {
		return errors.New("no mode class name found")
	}

	name := strings.TrimSpace(strings.ReplaceAll(result[0], "'", ""))
	gameMode := &models.Mode{}
	if err := p.db.Conn.Where(&models.Mode{Class: name}).First(gameMode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			gameMode.Name = name
			gameMode.Class = name
			if err := p.db.Conn.Create(gameMode).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	p.logger.Info(fmt.Sprintf("Game mode is now: %s", gameMode.Name))
	p.CurrentGame.Mode = gameMode

	return nil
}
