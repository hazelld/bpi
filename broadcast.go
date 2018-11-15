package bpi

// The broadcast type in their json is pretty large and ugly, so I am taking some liberties
// with how I am structuring it. It is mainly peeling apart some of the depth of the JSON
// to make it more flat.
//
// The Broadcast > Audio portion of the JSON tree moved into other places. The streams are put
// into the top level Broadcast.AudioStreams, and the Broadcasts are placed into the
// Broadcast.Broadcasters array (explained below on Broadcaster struct)
//
type Broadcast struct {
	Broadcasters []Broadcaster
	Video        Video
	AudioStreams []Stream
}

func LoadBroadcastFromScoreboard(scoreboard interface{}) Broadcast {
	broadcast := Broadcast{}

	// Load the audio & video broadcasters
	broadcast.Broadcasters = LoadBroadcastersFromScoreboardVideo(scoreboard)
	audio_bcasters := LoadBroadcastersFromScoreboardAudio(scoreboard)
	broadcast.Broadcasters = append(broadcast.Broadcasters, audio_bcasters...)

	// Load the video & audio streams + info
	broadcast.Video = LoadVideoFromScoreboard(scoreboard)
	broadcast.AudioStreams = LoadAudioStreamsFromScoreboard(scoreboard)

	return broadcast
}
