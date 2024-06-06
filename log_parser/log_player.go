package log_parser

import (
	"errors"
	"strconv"
	"strings"
)

type Player struct {
	SteamID string
	Name    string
	Team    uint
}

func (p *LogParser) parsePlayers(str string, simple bool) (*[]Player, error) {

	tempPlayers := strings.Split(str, " + ")
	if len(tempPlayers) == 0 {
		return nil, errors.New("no players found")
	}
	//fmt.Println(tempPlayers)
	players := make([]Player, 0)
	//fmt.Printf("count: %d\n", len(tempPlayers))
	for _, tempPlayer := range tempPlayers {
		var player *Player
		var err error
		//fmt.Println(tempPlayer)
		if simple {

			player, err = p.parsePlayerSimple(tempPlayer)
		} else {
			player, err = p.parsePlayer(tempPlayer)
		}

		if err != nil {
			return nil, err
		}
		players = append(players, *player)
	}

	if len(players) == 0 {
		return nil, errors.New("no players found")
	}

	return &players, nil
}

func (p *LogParser) parsePlayer(playerString string) (*Player, error) {
	//fmt.Println(playerString)
	idx := strings.Index(playerString, "[")
	if idx == -1 {
		return nil, errors.New("invalid player record")
	}
	name := strings.TrimSpace(playerString[:idx])
	idAndTeam := playerString[idx:]

	player_ids := p.regexps["player_get_id"].FindStringSubmatch(idAndTeam)
	if len(player_ids) == 0 {
		return nil, errors.New("no player found in record")
	}
	player_id := player_ids[0]
	player_id = strings.TrimSpace(strings.Replace(player_id, "[", "", 1))

	player_teams := p.regexps["player_get_team"].FindStringSubmatch(idAndTeam)
	if len(player_teams) == 0 {
		return nil, errors.New("no player team found in record")
	}

	player_team := player_teams[0]
	player_team = strings.TrimSpace(strings.Replace(strings.Replace(player_team, "team", "", 1), "]", "", 1))
	team, err := strconv.Atoi(player_team)
	if err != nil {
		return nil, errors.New("invalid team number")
	}

	player := Player{
		Name:    name,
		SteamID: player_id,
		Team:    uint(team),
	}

	return &player, nil
}

func (p *LogParser) parsePlayerSimple(playerString string) (*Player, error) {

	idx := strings.Index(playerString, "[")
	//fmt.Println(playerString)
	name := playerString[:idx]
	player_id := strings.TrimSpace(strings.Replace(strings.Replace(playerString[idx:], "]", "", 1), "[", "", 1))
	player := Player{
		Name:    name,
		SteamID: player_id,
		Team:    0,
	}

	return &player, nil
}
