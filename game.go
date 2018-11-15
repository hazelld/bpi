package bpi

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Game struct {
	League            string
	GameID            string
	SeasonStageID     int
	GameURLCode       string
	StatusNum         int
	ExtendedStatusNum int
	StartTimeUTC      string
	StartDateEastern  string
	StartTimeEastern  string
	IsBuzzerBeater    bool
	Tags              []string

	Period       Period
	Nugget       Nugget
	HTeam        TeamScoreboard
	VTeam        TeamScoreboard
	Video        Video
	Broadcasters []Broadcaster
}

func Games(year string) ([]Game, error) {
	all_games := []Game{}
	raw_json, err := MakeRequest(fmt.Sprintf("/prod/v2/%s/schedule.json", year))
	if err != nil {
		return all_games, err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(raw_json), &result)

	league_schedule := unwrapPath(result, []string{"league"})
	games_raw := unwrapMap(league_schedule, func(league_name string, d interface{}) interface{} {
		games := unwrapArray(d, func(i int, game_raw interface{}) interface{} {
			game := Game{}
			mapstructure.WeakDecode(game_raw, &game)
			game.League = league_name
			game.Video = LoadVideoFromSchedule(game_raw)
			game.Broadcasters = LoadBroadcastersFromSchedule(game_raw)
			return game
		})
		return games
	})

	for _, games_array := range games_raw {
		for _, game_raw := range games_array.([]interface{}) {
			game := game_raw.(Game)
			fmt.Println(game)
		}
	}
	return all_games, nil
}
