package bpi

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strconv"
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
	HasGamePDF            bool
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

//
type Arena struct {
	Name       string
	IsDomestic bool
	City       string
	StateAbbr  string
	Country    string
}

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

type GameDuration struct {
	Hours   int
	Minutes int // Making the assumption here that NBA doesn't use a decimal or something
}

type Period struct {
	Current       int
	Type          int
	MaxRegular    int // Incase they ever change the amount of quarters in a game?????
	IsHalftime    bool
	IsEndOfPeriod bool
}

type TeamScoreboard struct {
	TeamID     string
	TriCode    string
	Win        int
	Loss       int
	SeriesWin  int
	SeriesLoss int
	Score      int
	LineScore  []int
}

// The broadcast type in their json is pretty large and ugly, so I am taking some liberties
// with how I am structuring it. It is mainly peeling apart some of the depth of the JSON
// to make it more flat.
//
// The Broadcast > Audio portion of the JSON tree moved into other places. The streams are put
// into the top level Broadcast.AudioStreams, and the Broadcasts are placed into the
// Broadcast.Broadcasters array (explained below on Broadcaster struct)
//
// TODO: Is this the best way to do this? Should this be maybe broken into it's own part?
type Broadcast struct {
	Broadcasters []Broadcaster
	Video        Video
	AudioStreams []Stream
}

// Instead of having the following chain of structures (if following nba json exactly):
// - Broadcast > Broadcasters > { national, vTeam, etc.. } > { ShortName, LongName }
//
// The level with the key being the location of the broadcaster (national, etc..) is
// being added as a field in the Broadcaster type defined below. The new field is
// "Location".
//
// Note that the Audio Broadcasters that are located at:
// - Broadcast > Audio > { national, vTeam, hTeam } > Broadcasters
// are folded into this structure, so "Broadcaster" represents both the "broadcast > broadcasters"
// part of the JSON and the "broadcast > audio > brodcasters" part.
type Broadcaster struct {
	Type      string
	Location  string
	ShortName string
	LongName  string
}

type Video struct {
	RegionalBlackoutCodes string
	CanPurchase           bool
	IsLeaguePass          bool
	IsNationalBlackout    bool
	IsTntOt               bool // TNT-OT refered as TntOt in all places
	IsVR                  bool
	TntOtIsOnAir          bool
	IsMagicLeap           bool
	IsOculusVenues        bool
	Streams               []Stream
}

type Stream struct {
	StreamType            string // 'national', 'vTeam', 'hTeam' etc
	StreamFormat          string // Either 'video' or 'audio'
	IsOnAir               bool
	DoesArchiveExist      bool
	IsArchiveAvailToWatch bool
	StreamID              string
	Duration              int // Is this an int?
}

func Scoreboards(date string) ([]Scoreboard, error) {
	scoreboards := []Scoreboard{}
	raw_json, err := MakeRequest(fmt.Sprintf("/prod/v1/%s/scoreboard.json", date))
	if err != nil {
		return scoreboards, err
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(raw_json), &result)
	scoreboards_json := result["games"].([]interface{})

	// TODO: Find way to make this not so big & hacky
	for _, scoreboard_data := range scoreboards_json {

		// Fill out most of the struct with mapstructure - Then manually fill trickier parts
		scoreboard := Scoreboard{}
		mapstructure.WeakDecode(scoreboard_data, &scoreboard)

		// Flatten Linescore through hax (A guide to abusing type conversions in horrible ways)
		linescore_array_visitor := ((scoreboard_data.(map[string]interface{}))["vTeam"].(map[string]interface{}))["linescore"]
		linescore_array_home := ((scoreboard_data.(map[string]interface{}))["hTeam"].(map[string]interface{}))["linescore"]
		for i, ls := range linescore_array_visitor.([]interface{}) {
			scoreboard.VTeam.LineScore[i], _ = strconv.Atoi((ls.(map[string]interface{})["score"]).(string))
		}
		for i, ls := range linescore_array_home.([]interface{}) {
			scoreboard.HTeam.LineScore[i], _ = strconv.Atoi((ls.(map[string]interface{})["score"]).(string))
		}

		// Deal with the Broadcast shitshow
		broadcast := ((scoreboard_data.(map[string]interface{}))["watch"]).(map[string]interface{})["broadcast"]
		broadcasters := ((broadcast.(map[string]interface{}))["broadcasters"]).(map[string]interface{})

		// Note, it appears that under watch > broadcast > broadcasters > {national, vTeam, etc..}
		// that there _could_ be an array of these broadcasters (however I have yet to see any with
		// more than one element). I am going to make multiple entries in Broadcasters[] for each one
		// just incase
		for location, data := range broadcasters {

			bcaster_array := data.([]interface{})
			for _, bcaster_names := range bcaster_array {
				bcaster := Broadcaster{}
				mapstructure.WeakDecode(bcaster_names, &bcaster)
				bcaster.Location = location
				bcaster.Type = "video"
				scoreboard.Broadcast.Broadcasters = append(scoreboard.Broadcast.Broadcasters, bcaster)
			}
		}

		// Handle video portion - Have to flatten video > streams part of tree manually
		video := ((broadcast.(map[string]interface{}))["video"]).(map[string]interface{})
		mapstructure.WeakDecode(video, &scoreboard.Broadcast.Video)
		for _, stream_data := range video["streams"].([]interface{}) {
			stream := Stream{}
			mapstructure.WeakDecode(stream_data, &stream)
			stream.StreamFormat = "video"
			scoreboard.Broadcast.Video.Streams = append(scoreboard.Broadcast.Video.Streams, stream)
		}

		// Get the Audio Stream portions
		audio := ((broadcast.(map[string]interface{}))["audio"]).(map[string]interface{})
		for provider, btype := range audio {
			streams := btype.(map[string]interface{})["streams"].([]interface{})

			for _, stream := range streams {
				tmp_stream := Stream{}
				mapstructure.WeakDecode(stream, &tmp_stream)
				tmp_stream.StreamFormat = "audio"
				tmp_stream.StreamType = provider
				scoreboard.Broadcast.AudioStreams = append(scoreboard.Broadcast.AudioStreams, tmp_stream)
			}

			audio_bcasters := btype.(map[string]interface{})["broadcasters"].([]interface{})
			for _, audio_bcast_data := range audio_bcasters {
				tmp_bcaster := Broadcaster{}
				mapstructure.WeakDecode(audio_bcast_data, &tmp_bcaster)
				tmp_bcaster.Type = "audio"
				tmp_bcaster.Location = provider
				scoreboard.Broadcast.Broadcasters = append(scoreboard.Broadcast.Broadcasters, tmp_bcaster)
			}
		}

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
