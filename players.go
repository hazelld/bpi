package bpi

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Player struct {
	FirstName    string
	LastName     string
	PersonID     string
	TeamID       string
	Jersey       int
	IsActive     bool
	Pos          string
	HeightFeet   int
	HeightInches int

	// Chose to store these as strings instead of float just to avoid rounding issues
	HeightMeters    string
	WeightPounds    string
	WeightKilograms string

	// Store the date as a string, up to consumer to convert
	DateOfBirthUTC  string
	NBADebutYear    string
	YearsPro        string
	CollegeName     string
	LastAffiliation string
	Country         string

	// Which league the player is in. Values are: {standard, africa, sacramento, vegas, utah}
	League string

	// Nested structs
	Teams []PlayerTeamData
	Draft PlayerDraftData
}

// This struct holds the information for where the player played between SeasonStart
// and SeasonEnd years.
type PlayerTeamData struct {
	TeamID      string
	SeasonStart string
	SeasonEnd   string
}

// This struct holds the information for where the player was drafted.
type PlayerDraftData struct {
	TeamID     string
	PickNum    string
	RoundNum   string
	SeasonYear string
}

// Load all the players for a given year. Note this will load the NBA players,
// summer league players, and the African league players.
func Players(year string) ([]Player, error) {
	all_players := []Player{}
	raw_json, err := MakeRequest(fmt.Sprintf("/prod/v1/%s/players.json", year))
	if err != nil {
		return all_players, err
	}

	// Unpeel top level (league)
	var result map[string]interface{}
	json.Unmarshal([]byte(raw_json), &result)

	player_json := result["league"].(map[string]interface{})

	for league_type, players := range player_json {
		for _, player_data := range players.([]interface{}) {
			player := Player{}
			mapstructure.WeakDecode(player_data, &player)
			player.League = league_type
			all_players = append(all_players, player)
		}
	}
	return all_players, nil
}

//
func FilterPlayers(players []Player, test func(Player) bool) (filtered_players []Player) {
	for _, player := range players {
		if test(player) {
			filtered_players = append(filtered_players, player)
		}
	}
	return
}
