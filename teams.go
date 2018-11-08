package bpi

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Team struct {
	IsNBAFranchise bool   `json:"isNBAFranchise"`
	IsAllStar      bool   `json:"isAllStar"`
	City           string `json:"city"`
	AltCityName    string `json:"altCityName"`
	FullName       string `json:"fullName"`
	TriCode        string `json:"tricode"`
	TeamID         string `json:"teamId"`
	NickName       string `json:"nickname"`
	URLName        string `json:"urlName"`
	ConfName       string `json:"confName"`
	DivName        string `json:"divName"`
}

func Teams(year string) ([]Team, error) {
	teams := []Team{}
	raw_json, err := MakeRequest(fmt.Sprintf("/prod/v2/%s/teams.json", year))

	if err != nil {
		return teams, err
	}

	// Need to go a few levels deep to get the proper part of the json tree, want
	// to unpeel the layers: league > standard. Pretty hacky but it works, as we
	// can map the unpeeled []interface{} to an array of Team structures using
	// mapstructure
	var result map[string]interface{}
	json.Unmarshal([]byte(raw_json), &result)

	teams_json := result["league"].(map[string]interface{})
	for _, dat := range teams_json["standard"].([]interface{}) {
		team := Team{}
		mapstructure.WeakDecode(dat, &team)
		teams = append(teams, team)
	}
	return teams, nil
}

// Get just the NBA teams
func NBATeams(year string) ([]Team, error) {
	teams, err := Teams(year)

	if err != nil {
		return teams, err
	}

	return FilterTeams(teams, func(t Team) bool {
		return t.IsNBAFranchise == true
	}), nil
}

func FilterTeams(teams []Team, test func(Team) bool) (filtered_teams []Team) {
	for _, team := range teams {
		if test(team) {
			filtered_teams = append(filtered_teams, team)
		}
	}
	return
}
