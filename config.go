package main

import (
	"strings"

	"github.com/imkira/go-interpol"
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

func MergeMaps(maps ...map[string]string) (result map[string]string) {
	result = make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			result[strings.ToLower(k)] = v
		}
	}

	return result
}

func FormatMap(status map[string]string) map[string]string {
	var values = map[string]string{}

	constants := []string{
		"volume",
		"repeat",
		"random",
		"single",
		"playlistlength",
		"consume",
		"audio",
		"bitrate",
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
		"artists",
		"albums",
		"songs",
	}

	for _, s := range constants {
		values[strings.ToLower(s)] = status[strings.Title(s)]
	}

	return values
}

func Formatted(s string, m map[string]string) string {
    formatted, err := interpol.WithMap(s, m)
    if err != nil {
        panic(err)
    }

    return formatted
}
