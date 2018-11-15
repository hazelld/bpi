package bpi

import (
	//	"fmt"
	"github.com/mitchellh/mapstructure"
)

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

func LoadBroadcastersFromScoreboardVideo(scoreboard interface{}) []Broadcaster {
	broadcasters := []Broadcaster{}
	bcaster_data := unwrapPath(scoreboard, []string{"watch", "broadcast", "broadcasters"})

	parsed_bcaster := unwrapMap(bcaster_data, func(loc string, data interface{}) interface{} {
		bcasters := unwrapArray(data, func(i int, bcast interface{}) interface{} {
			bcaster := Broadcaster{}
			mapstructure.WeakDecode(bcast, &bcaster)
			bcaster.Location = loc
			bcaster.Type = "video"
			return bcaster
		})
		return bcasters
	})

	// Since each parsed_bcaster[i] is an array of broadcasters, have to unpeel
	for _, bcaster_arr := range parsed_bcaster {
		for _, bcaster := range bcaster_arr.([]interface{}) {
			broadcasters = append(broadcasters, bcaster.(Broadcaster))
		}
	}
	return broadcasters
}

func LoadBroadcastersFromScoreboardAudio(scoreboard interface{}) []Broadcaster {
	broadcasters := []Broadcaster{}
	audio_data := unwrapPath(scoreboard, []string{"watch", "broadcast", "audio"})

	parsed_bcaster := unwrapMap(audio_data, func(loc string, data interface{}) interface{} {
		bcaster_data := unwrapPath(data, []string{"broadcasters"})
		bcasters := unwrapArray(bcaster_data, func(i int, bcast interface{}) interface{} {
			bcaster := Broadcaster{}
			mapstructure.WeakDecode(bcast, &bcaster)
			bcaster.Location = loc
			bcaster.Type = "audio"
			return bcaster
		})
		return bcasters
	})

	// Since each parsed_bcaster[i] is an array of broadcasters, have to unpeel
	for _, bcaster_arr := range parsed_bcaster {
		for _, bcaster := range bcaster_arr.([]interface{}) {
			broadcasters = append(broadcasters, bcaster.(Broadcaster))
		}
	}
	return broadcasters
}

// This loads the broadcasters from the Game portion of the schedule JSON. This part is
// really weird, and is a pain in the ass to parse. The broadcasters are top-level parts
// of this part of the JSON, and have a path like: "national" : { "broadcasters" : [...] }
// To be safe here, we will unwrap the top level map, and if the item is a
// map[string]interface{}, we check for the "broadcasters" key. This should be less hacky
// and more robust than checking for the key of the broadcaster, since those may change
// over time. It's a massive pain in the ass, since the other places we load the
// broadcasters the _only_ keys in that part of the tree are the broadcaster locations,
// while here the regular Video data is also there.
//
// Note for some dumb reason the key "national" actually has the sub-tree "broadcasters"
// that has the info, BUT the "canadian" and "spanish_national" DONT and are just an
// array. Noooo idea why.
func LoadBroadcastersFromSchedule(game_data interface{}) []Broadcaster {
	video_data := unwrapPath(game_data, []string{"watch", "broadcast", "video"})

	parsed_bcaster := unwrapMap(video_data, func(key string, d interface{}) interface{} {
		bcs, ok := d.(map[string]interface{})
		var bcasters_raw interface{}

		// If it was a map[string]interface{} unwrap the "broadcasters" part
		if !ok {
			bcasters_raw, ok = d.([]interface{})
			if !ok {
				return nil
			}
		} else {
			bcasters_raw = unwrapPath(bcs, []string{"broadcasters"})
		}

		bcasters_arr := unwrapArray(bcasters_raw, func(i int, bc interface{}) interface{} {
			bcaster := Broadcaster{}
			mapstructure.WeakDecode(bc, &bcaster)
			bcaster.Type = "video"
			bcaster.Location = key
			return bcaster
		})
		return bcasters_arr
	})

	broadcasters := []Broadcaster{}
	// Since each parsed_bcaster[i] is an array of broadcasters, have to unpeel
	for _, bcaster_arr := range parsed_bcaster {

		// nils are from the keys in the video json that weren't broadcasters
		if bcaster_arr == nil {
			continue
		}

		for _, bcaster := range bcaster_arr.([]interface{}) {
			broadcasters = append(broadcasters, bcaster.(Broadcaster))
		}
	}
	return broadcasters
}
