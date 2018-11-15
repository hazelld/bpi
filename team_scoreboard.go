package bpi

import (
//    "fmt"
    "strconv"
    "github.com/mitchellh/mapstructure"
)

type TeamScoreboard struct {
    GameLocation string // vTeam or hTeam
	TeamID     string
	TriCode    string
	Win        int
	Loss       int
	SeriesWin  int
	SeriesLoss int
	Score      int
	LineScore  []int
}

// Load the TeamScoreboard struct from the context of the Scoreboard type.
//
// The @team param should have either: "vTeam" or "hTeam" depending on which
// team's data is loaded.
func LoadTeamScoreboardFromScoreboard(scoreboard interface{}, team string) TeamScoreboard {
    team_s := TeamScoreboard{}
    team_data := unwrapPath(scoreboard, []string{team})
    mapstructure.WeakDecode(team_data, &team_s)

    linescore := unwrapPath(scoreboard, []string{team, "linescore"})
    linescore_collection := unwrapArray(linescore, func(i int, d interface{}) interface{} {
        score_data := unwrapPath(d, []string{"score"})
        return score_data
    })

    // Convert each of the scores from interface{} to string, then to int
    for i, ls := range linescore_collection {
        team_s.LineScore[i], _ = strconv.Atoi(ls.(string))
    }
    team_s.GameLocation = team
    return team_s
}
