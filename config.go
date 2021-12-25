package main

import (
	"fmt"
	"strings"

	"github.com/imkira/go-interpol"
)

// define default config
type Format struct {
	Details, State, LargeText string
	Remaining                 bool
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
		LargeText: "{album}",
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
		"album",
		"albumartist",
		"albums",
		"artist",
		"artists",
		"audio",
		"bitrate",
		"composer",
		"conductor",
		"consume",
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
		"playlistlength",
		"random",
		"repeat",
		"single",
		"songs",
		"title",
		"track",
		"volume",
		"work",
	}

	for _, s := range constants {
		values[strings.ToLower(s)] = status[strings.Title(s)]
	}

	return values
}

func Formatted(s string, m map[string]string) string {
	formatted, err := interpol.WithMap(s, m)
	if err != nil {
		fmt.Println(err)
	}

	return formatted
}
