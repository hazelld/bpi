package bpi

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"time"
)

type Play struct {
	// These aren't in the JSON, but provide context to when the play happened
	Date   time.Time
	GameID string
	Period int

	// JSON Fields
	Clock                string
	EventMsgType         int // Can this be int?
	Description          string
	PersonID             int
	TeamID               int
	VTeamScore           int
	HTeamScore           int
	IsScoreChange        bool
	IsVideoAvailable     bool
	FormattedDescription string
}

func PlaysByGame(date string, gameid string) []Play {
	all_plays := []Play{}

	for _, i := range []int{1, 2, 3, 4} {
		period_plays, err := PlaysByGameAndPeriod(date, gameid, i)

		if err == nil {
			all_plays = append(all_plays, period_plays...)
		}
	}

	return all_plays
}

func PlaysByGameAndPeriod(date string, gameid string, period int) ([]Play, error) {
	plays := []Play{}
	endpoint := fmt.Sprintf("/prod/v1/%s/%s_pbp_%d.json", date, gameid, period)
	raw_json, err := MakeRequest(endpoint)
	if err != nil {
		return plays, err
	}

	gamedate, err := time.Parse("20060102", date)
	if err != nil {
		return plays, err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(raw_json), &result)
	plays_raw := unwrapPath(result, []string{"plays"})
	all_plays_raw := unwrapArray(plays_raw, func(i int, d interface{}) interface{} {
		play := Play{
			Date:   gamedate,
			GameID: gameid,
			Period: period,
		}
		mapstructure.WeakDecode(d, &play)
		return play
	})

	for _, i := range all_plays_raw {
		plays = append(plays, i.(Play))
	}
	return plays, nil
}
