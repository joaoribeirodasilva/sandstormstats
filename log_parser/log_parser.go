package log_parser

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joaoribeirodasilva/sandstormstats/db"
	"github.com/joaoribeirodasilva/sandstormstats/models"
	"gorm.io/gorm"
)

type LogParser struct {
	Servers     []models.Server
	db          *db.Db
	logger      *slog.Logger
	regexps     map[string]*regexp.Regexp
	CurrentTime time.Time
	CurrentGame LogGame
}

type LogGame struct {
	Server *models.Server
	Game   *models.Game
	Round  *models.Round
	Map    *models.Map
	Mode   *models.Mode
}

type LogFile struct {
	Name string
	Size int64
	Line int
}

func New(db *db.Db, logger *slog.Logger) *LogParser {

	p := new(LogParser)
	p.logger = logger
	p.db = db
	p.regexps = make(map[string]*regexp.Regexp, 0)
	return p
}

func (p *LogParser) Parse() {

	if err := p.registerRegExprs(); err != nil {
		return
	}

	p.logger.Info("starting log parser")
	if err := p.getServers(); err != nil {
		p.logger.Error(err.Error())
		return
	}

	for _, server := range p.Servers {

		p.CurrentGame.Game = nil
		p.CurrentGame.Round = nil
		p.CurrentGame.Map = nil
		p.CurrentGame.Mode = nil

		p.CurrentGame.Server = &server
		if err := p.getLastServerGame(p.CurrentGame.Server); err != nil {
			p.logger.Error(fmt.Sprintf("error getting last server %s game from database", p.CurrentGame.Server.Key))
			os.Exit(1)
		}
		if err := p.parseServer(&server); err != nil {
			p.logger.Error(err.Error())
			continue
		}
	}
}

func (p *LogParser) parseServer(server *models.Server) error {

	var err error
	var logs *[]LogFile

	logs, err = p.getFiles(server)
	if err != nil {
		return err
	}

	if len(*logs) == 0 {
		p.logger.Debug(fmt.Sprintf("no logs found for server %s", server.Key))
		return nil
	}

	for _, log := range *logs {
		server_log := &models.ServerLog{
			ServerID:  server.ID,
			File:      log.Name,
			StartTime: time.Now().UTC(),
			EndTime:   nil,
			Success:   0,
			Message:   "",
		}

		/* if err = p.db.Conn.Create(server_log).Error; err != nil {
			p.logger.Error(err.Error())
			return err
		} */

		logPath := fmt.Sprintf("%s/%s", server.LogDir, log.Name)
		if err := p.parseFile(logPath, server); err != nil {
			now := time.Now().UTC()
			server_log.EndTime = &now
			server_log.Message = err.Error()

			server.LastLog = &log.Name
			server.LastLogSize = uint64(log.Size)
			/* if err = p.db.Conn.Save(server).Error; err != nil {
				p.logger.Error(err.Error())
				return err
			} */
		} else {
			now := time.Now().UTC()
			server_log.EndTime = &now
			server_log.Success = 1
		}

		/* if err = p.db.Conn.Save(server_log).Error; err != nil {
			p.logger.Error(err.Error())
			return err
		} */
	}

	return nil
}

func (p *LogParser) parseFile(path string, server *models.Server) error {

	p.logger.Info(fmt.Sprintf("parsing log file %s for server %s", path, server.Key))

	file, err := os.Open(path)
	if err != nil {
		p.logger.Error(fmt.Sprintf("error opening file: %v\n", err))
		return err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		str := fileScanner.Text()
		server.LastLine++
		//fmt.Printf("\rLine: %d\n", server.LastLine)
		//fmt.Printf("[%d]: %s\n", server.LastLine, str)
		dateTime := p.HasDateTime(str)
		if dateTime == "" || !p.HasValidEntry(str) {
			continue
		}
		p.CurrentTime, err = time.Parse("[2006.02.01-15.04.05.000]", dateTime)
		if err != nil {
			//p.logger.Error(err.Error())
			p.logger.Error(fmt.Sprintf("invalid log timestamp %s on file %s, Line: %d", dateTime, path, server.LastLine))
			return err
		}

		if err := p.kill(str); err != nil {
			p.logger.Error(fmt.Sprintf("%s on file %s, Line: %d", err.Error(), path, server.LastLine))
			return err
		}

		if err := p.objective(str); err != nil {
			p.logger.Error(fmt.Sprintf("%s on file %s, Line: %d", err.Error(), path, server.LastLine))
			return err
		}

		if err := p.round(str); err != nil {
			p.logger.Error(fmt.Sprintf("%s on file %s, Line: %d", err.Error(), path, server.LastLine))
			return err
		}

		if err := p.scenario(str); err != nil {
			p.logger.Error(fmt.Sprintf("%s on file %s, Line: %d", err.Error(), path, server.LastLine))
			return err
		}

		if err := p.mode(str); err != nil {
			p.logger.Error(fmt.Sprintf("%s on file %s, Line: %d", err.Error(), path, server.LastLine))
			return err
		}
	}

	return nil
}

func (p *LogParser) HasDateTime(str string) string {

	result := p.regexps["has_date"].FindStringSubmatch(str)
	if len(result) == 0 {
		return ""
	}
	return strings.Replace(result[0], ":", ".", 1)
}

func (p *LogParser) HasValidEntry(str string) bool {
	result := p.regexps["has_valid_entry"].FindStringSubmatch(str)
	return len(result) != 0
}

func (p *LogParser) getFiles(server *models.Server) (*[]LogFile, error) {

	p.logger.Debug(fmt.Sprintf("finding log files for server %s", server.Key))

	logFiles := make([]LogFile, 0)
	items, err := os.ReadDir(server.LogDir)
	if err != nil {
		p.logger.Error(err.Error())
		return nil, err
	}

	lastLog := ""
	if server.LastLog != nil {
		lastLog = *server.LastLog
	}

	p.logger.Debug(fmt.Sprintf("last log for server %s is %s", server.Key, lastLog))

	for _, item := range items {
		if !item.IsDir() {
			logFileInfo, err := item.Info()
			if err != nil {
				p.logger.Error(err.Error())
				return nil, err
			}

			if strings.HasPrefix(item.Name(), server.Key) && item.Name() != (server.Key+".log") && item.Name() != (server.Key+"-CRC.log") && item.Name() >= lastLog {

				if item.Name() == lastLog && logFileInfo.Size() == int64(server.LastLogSize) {
					p.logger.Debug(fmt.Sprintf("server %s log file %s already fully processed", server.Key, lastLog))
					continue
				}

				logFile := LogFile{
					Name: item.Name(),
					Size: logFileInfo.Size(),
					Line: 0,
				}

				p.logger.Debug(fmt.Sprintf("adding log file %s for server %s to the parse queue", logFile.Name, server.Key))
				//fmt.Printf("Log: %+v\n", logFile)
				logFiles = append(logFiles, logFile)

			} else {
				p.logger.Debug(fmt.Sprintf("server %s log file %s already processed", server.Key, lastLog))
			}
		}
	}

	return &logFiles, nil
}

func (p *LogParser) getServers() error {

	if err := p.db.Conn.Model(&models.Server{}).Find(&p.Servers).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	//fmt.Printf("Servers: %+v\n", p.Servers)
	return nil
}

func (p *LogParser) getLastServerGame(server *models.Server) error {
	game := &models.Game{}
	if err := p.db.Conn.Preload("Map").Where(&models.Game{ServerID: server.ID}).Order("id desc").First(game).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.CurrentGame.Game = nil
			return nil
		}
		return err
	}
	p.CurrentGame.Game = game
	return nil
}

func (p *LogParser) registerRegExprs() error {

	var err error
	p.regexps["is_capture"], err = regexp.Compile(`LogGameplayEvents: Display: Objective \d was .* by .*`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("is_capture regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["is_destroy"], err = regexp.Compile(`LogGameplayEvents: Display: Objective \d owned .* by .*`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("is_destroy regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["is_scenario"], err = regexp.Compile(`LogScenario: Display: Loading scenario '.*`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("is_scenario regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["scenario_file"], err = regexp.Compile(`'.*'`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("scenario_file regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["is_mode"], err = regexp.Compile(`LogLoad: Game class is '.*`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("is_mode regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["mode_class"], err = regexp.Compile(`'.*'`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("mode_class regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["is_start_round"], err = regexp.Compile(`LogGameplayEvents: Display: Round \d* started`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("is_start_round regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["round_number"], err = regexp.Compile(`Round \d*`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("round_number regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["is_end_round"], err = regexp.Compile(`LogGameplayEvents: Display: Round \d* Over:`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("is_end_round regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["round_end_team"], err = regexp.Compile(`Team \d*`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("round_end_team regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["is_kill"], err = regexp.Compile(`LogGameplayEvents: Display: .* killed`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("is_kill regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["has_date"], err = regexp.Compile(`^\[\d{4}\.\d{2}\.\d{2}\-\d{2}\.\d{2}\.\d{2}\:\d{3}\]`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("has_date regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["has_valid_entry"], err = regexp.Compile(`(LogGameplayEvents)|(LogGameMode)|(LogLoad)|(LogScenario)`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("has_valid_entry regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["kill_get_players"], err = regexp.Compile(`.*\[\d*, team \d{1}\]`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("has_valid_entry regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["kill_get_weapon"], err = regexp.Compile(`with .*$`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("has_valid_entry regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["player_get_id"], err = regexp.Compile(`\[(\d*)|(INVALID\.*)`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("player_get_id regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["player_get_id2"], err = regexp.Compile(`\[\d{17}\]`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("player_get_id regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["player_get_team"], err = regexp.Compile(`team \d\]`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("player_get_team regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["objective_get_num"], err = regexp.Compile(`Objective \d`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("player_get_team regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["objective_get_num"], err = regexp.Compile(`Objective \d`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("objective_get_num regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["objective_get_team"], err = regexp.Compile(`for team \d`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("objective_get_team regular expression error: %s", err.Error()))
		return err
	}

	p.regexps["objective_get_players"], err = regexp.Compile(`by .*\]`)
	if err != nil {
		p.logger.Error(fmt.Sprintf("objective_get_players regular expression error: %s", err.Error()))
		return err
	}

	return nil
}
