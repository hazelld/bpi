package bpi

import (
	//	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Stream struct {
	StreamType            string // 'national', 'vTeam', 'hTeam' etc
	StreamFormat          string // Either 'video' or 'audio'
	Language              string
	IsOnAir               bool
	DoesArchiveExist      bool
	IsArchiveAvailToWatch bool
	StreamID              string
	Duration              int // Is this an int?
}

func LoadVideoStreamsFromScoreboard(scoreboard interface{}) []Stream {
	stream_arr := unwrapPath(scoreboard, []string{"watch", "broadcast", "video", "streams"})
	streams_data := unwrapArray(stream_arr, func(i int, d interface{}) interface{} {
		stream := Stream{}
		mapstructure.WeakDecode(d, &stream)
		stream.StreamFormat = "Video"
		return stream
	})

	streams := []Stream{}
	for _, data := range streams_data {
		streams = append(streams, data.(Stream))
	}
	return streams
}

func LoadAudioStreamsFromScoreboard(scoreboard interface{}) []Stream {
	stream_arr := unwrapPath(scoreboard, []string{"watch", "broadcast", "audio"})
	streams_data := unwrapMap(stream_arr, func(stype string, d interface{}) interface{} {
		audio_streams_data := unwrapPath(d, []string{"streams"})
		return unwrapArray(audio_streams_data, func(i int, d interface{}) interface{} {
			stream := Stream{}
			mapstructure.WeakDecode(d, &stream)
			stream.StreamFormat = "Audio"
			stream.StreamType = stype
			return stream
		})
	})

	// Have to unpack an array of arrays of audio streams (unwrapMap then unwrapArray call)
	streams := []Stream{}
	for _, data := range streams_data {
		for _, stream_data := range data.([]interface{}) {
			streams = append(streams, stream_data.(Stream))
		}
	}
	return streams
}
