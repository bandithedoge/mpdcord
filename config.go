package main

import (
	"fmt"
	"strings"

	"github.com/imkira/go-interpol"
)

// define default config
type Format struct {
	Details, State, LargeText, SmallText string
	Remaining                            bool
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
		SmallText: "{state}",
		Remaining: false,
	},
}

func MergeMaps(maps ...map[string]string) map[string]string {
	result := map[string]string{}

	for _, m := range maps {
		for k, v := range m {
			result[strings.ToLower(k)] = v
		}
	}

	return result
}

func Formatted(s string, m map[string]string) string {
	formatted, err := interpol.WithMap(s, m)
	if err != nil {
		fmt.Println(s, ": ", err)
	}

	return formatted
}
