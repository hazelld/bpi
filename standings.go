package bpi

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type TeamStanding struct {
	League                 string
	TeamID                 string
	Win                    int
	Loss                   int
	WinPct                 string
	WinPctV2               string
	LossPct                string
	LossPctV2              string
	GamesBehind            string // Could be float?
	DivGamesBehind         string
	ClinchedPlayoffsCode   string
	ClinchedPlayoffsCodeV2 string
	ConfRank               int
	ConfWin                int
	ConfLoss               int
	DivRank                int
	DivWin                 int
	DivLoss                int
	HomeWin                int
	HomeLoss               int
	AwayWin                int
	AwayLoss               int
	LastTenWin             int
	LastTenLoss            int
	Streak                 int
	IsWinStreak            bool
	TieBreakerPts          string
	SortKey                TeamSortKey
}

type TeamSortKey struct {
	DefaultOrder   int
	Nickname       int
	Win            int
	Loss           int
	WinPct         int
	GamesBehind    int
	ConfWinLoss    int
	DivWinLoss     int
	HomeWinLoss    int
	AwayWinLoss    int
	LastTenWinLoss int
	Streak         int
}

// Return the current standings
func CurrentStandings() ([]TeamStanding, error) {
	standings := []TeamStanding{}
	raw_json, err := MakeRequest(fmt.Sprintf("/prod/v1/current/standings_all.json"))
	if err != nil {
		return standings, err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(raw_json), &result)

	league := unwrapPath(result, []string{"league"})
	standings_raw := unwrapMap(league, func(league_name string, d interface{}) interface{} {
		teams := unwrapPath(d, []string{"teams"})
		league_standings := unwrapArray(teams, func(i int, standing_raw interface{}) interface{} {
			ts := TeamStanding{}
			mapstructure.WeakDecode(standing_raw, &ts)
			ts.League = league_name
			return ts
		})
		return league_standings
	})

	for _, league_standings := range standings_raw {
		for _, team_standing := range league_standings.([]interface{}) {
			standings = append(standings, team_standing.(TeamStanding))
		}
	}
	return standings, nil
}
