package main

import (
	"strings"

	"github.com/fhs/gompd/mpd"
)

// define default config
type Format struct {
	State, Details string
	Remaining      bool
}

type Config struct {
	ID                         int
	Address, Network, Password string
	Format                     Format
}

var DefaultConfig = Config{
	ID:       922175995828654100,
	Address:  "localhost:6600",
	Network:  "tcp",
	Password: "",
	Format: Format{
		Details:   "{title}",
		State:     "{artist}",
		Remaining: false,
	},
}

func FormatMap(status mpd.Attrs, song mpd.Attrs, stats mpd.Attrs) map[string]string {
	var values = map[string]string{}

	constants := map[string][]string{
		"status": {
			"volume",
			"repeat",
			"random",
			"single",
			"playlistlength",
			"consume",
			"audio",
			"bitrate",
		},
		"song": {
			"album",
			"artist",
			"albumartist",
			"composer",
			"conductor",
			"date",
			"disc",
			"ensemble",
			"genre",
			"grouping",
			"label",
			"location",
			"movement",
			"movementnumber",
			"originaldate",
			"performer",
			"title",
			"track",
			"work",
		},
		"stats": {
			"artists",
			"albums",
			"songs",
		},
	}

	for _, s := range constants["status"] {
        values[s] = status[strings.Title(s)]
	}
	for _, s := range constants["song"] {
        values[s] = song[strings.Title(s)]
	}
	for _, s := range constants["stats"] {
        values[s] = stats[strings.Title(s)]
	}

	return values
}
