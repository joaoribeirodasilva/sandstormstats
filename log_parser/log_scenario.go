package log_parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joaoribeirodasilva/sandstormstats/models"
	"gorm.io/gorm"
)

func (p *LogParser) scenario(logLine string) error {

	isS := p.regexps["is_scenario"].MatchString(logLine)
	if !isS {
		return nil
	}

	if err := p.parseScenario(logLine); err != nil {
		return err
	}

	return nil
}

func (p *LogParser) parseScenario(str string) error {

	result := p.regexps["mode_class"].FindStringSubmatch(str)
	if len(result) == 0 {
		return errors.New("no map name found")
	}

	name := strings.TrimSpace(strings.ReplaceAll(result[0], "'", ""))
	gameMap := &models.Map{}

	if err := p.db.Conn.Where(&models.Map{File: name}).First(gameMap).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			gameMap.Name = name
			gameMap.File = name
			if err := p.db.Conn.Create(gameMap).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	p.logger.Info(fmt.Sprintf("Game map is now: %s", gameMap.Name))
	p.CurrentGame.Map = gameMap

	return nil
}
