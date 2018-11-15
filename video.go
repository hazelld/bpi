package bpi

import (
	"github.com/mitchellh/mapstructure"
)

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

// Load the Video type from the Scoreboard. Note this will also load the streams that are
// contained within the Video.
func LoadVideoFromScoreboard(scoreboard interface{}) Video {
	video := Video{}
	video_data := unwrapPath(scoreboard, []string{"watch", "broadcast", "video"})
	mapstructure.WeakDecode(video_data, &video)

	// Get the streams
	video.Streams = LoadVideoStreamsFromScoreboard(scoreboard)
	return video
}

// Load the video from schedule.json file for the Game structure. Note that this requires
// the part of the JSON tree after "league" and "league_name" have been unwrapped. (ie the
// game portion of the JSON).
//
// Note this JSON is slightly different than the Video from scoreboard.json. It appears
// there is no streams associated, but there are fields that associate with a broadcaster?
// Such as: "national", "canadian". These are ignored in here, but _are_ a part of the
// Game structure.
func LoadVideoFromSchedule(game_data interface{}) Video {
	video := Video{}
	video_data := unwrapPath(game_data, []string{"watch", "broadcast", "video"})
	mapstructure.WeakDecode(video_data, &video)
	return video
}
