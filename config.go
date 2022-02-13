package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/imkira/go-interpol"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

func GetConfig() error {
	viper.SetConfigName("mpdcord.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(os.Getenv("XDG_CONFIG_HOME"))
	viper.AddConfigPath("$HOME/.config")

    err := viper.ReadInConfig()
	viper.Unmarshal(DefaultConfig)
	viper.OnConfigChange(func(e fsnotify.Event) {
		pterm.Info.Printfln("Config file changed: %s", e.Name)
	})
	viper.WatchConfig()

	return err
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
