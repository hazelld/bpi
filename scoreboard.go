package bpi

import (
	"encoding/json"
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
)

// Top level structure representing the /prod/v1/{date}/scoreboard.json file
// Other than the Broadcast part of this structure, it mostly maps 1-to-1 with
// the JSON file.
type Scoreboard struct {
	SeasonStageID         int
	SeasonYear            string
	GameID                string
	IsGameActivated       bool
	StatusNum             int
	ExtendedStatusNum     int
	StartTimeEastern      string
	StartTimeUTC          string
	StartDateEastern      string
	Clock                 string
	IsBuzzerBeater        bool
	IsPreviewArticleAvail bool
	IsRecapArticleAvail   bool
	HasGameBookPDF        bool
	IsStartTimeTBD        bool
	Attendance            string

	// Nested structures
	Arena        Arena
	Tickets      Tickets
	Nugget       Nugget // Wtf is this?
	GameDuration GameDuration
	Period       Period
	VTeam        TeamScoreboard
	HTeam        TeamScoreboard

	// Broadcast is long and ugly
	Broadcast Broadcast
}

// Note this will be loaded by the mapstructure call on the Scoreboard type. Arena's
// don't exist outside the context of the scoreboard.json as far as I can tell
type Arena struct {
	Name       string
	IsDomestic bool
	City       string
	StateAbbr  string
	Country    string
}

// Tickets are loaded through the Scoreboard type, and don't nee a special
// loading method.
type Tickets struct {
	MobileApp    string
	DesktopWeb   string
	MobileWeb    string
	LeagGameInfo string
	LeagTix      string
}

// This needs to be a sub-structure as the NBA JSON for some reason has it like this:
// { "nugget" : { "text": "nugget text here" } }
//
// So instead of manually correcting it, we will follow what their JSON does.
type Nugget struct {
	Text string
}

// This is loaded by the Scoreboard type, does not need a special loading method.
type GameDuration struct {
	Hours   int
	Minutes int // Making the assumption here that NBA doesn't use a decimal or something
}

// Loaded by the Scoreboard type.
type Period struct {
	Current       int
	Type          int
	MaxRegular    int // Incase they ever change the amount of quarters in a game?????
	IsHalftime    bool
	IsEndOfPeriod bool
}

// This function loads all the scoreboards for a given day
func Scoreboards(date string) ([]Scoreboard, error) {
	scoreboards := []Scoreboard{}
	raw_json, err := MakeRequest(fmt.Sprintf("/prod/v2/%s/scoreboard.json", date))
	if err != nil {
		return scoreboards, err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(raw_json), &result)
	scoreboards_json := result["games"].([]interface{})

	for _, scoreboard_data := range scoreboards_json {
		// Fill out most of the struct with mapstructure - Then manually fill trickier parts
		scoreboard := Scoreboard{}
		mapstructure.WeakDecode(scoreboard_data, &scoreboard)
		scoreboard.HTeam = LoadTeamScoreboardFromScoreboard(scoreboard_data, "hTeam")
		scoreboard.VTeam = LoadTeamScoreboardFromScoreboard(scoreboard_data, "vTeam")
		scoreboard.Broadcast = LoadBroadcastFromScoreboard(scoreboard_data)
		scoreboards = append(scoreboards, scoreboard)
	}
	return scoreboards, nil
}

func FilterScoreboards(scoreboards []Scoreboard, test func(Scoreboard) bool) (filtered_scoreboards []Scoreboard) {
	for _, scoreboard := range scoreboards {
		if test(scoreboard) {
			filtered_scoreboards = append(filtered_scoreboards, scoreboard)
		}
	}
	return
}
