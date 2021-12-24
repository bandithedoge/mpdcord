package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/akamensky/argparse"
	"github.com/dixonwille/wlog/v3"
	"github.com/fhs/gompd/mpd"
	"github.com/hugolgst/rich-go/client"
	"github.com/imdario/mergo"
	"github.com/imkira/go-interpol"
)

func main() {
	// setup logger
	var ui wlog.UI
	ui = wlog.New(os.Stdin, os.Stdout, os.Stdout)
	ui = wlog.AddPrefix("?", wlog.Cross, " ", "", "", "~", wlog.Check, "!", ui)
	ui = wlog.AddColor(wlog.Magenta, wlog.Red, wlog.Blue, wlog.BrightWhite, wlog.White, wlog.BrightMagenta, wlog.Cyan, wlog.Green, wlog.Yellow, ui)
	ui = wlog.AddConcurrent(ui)

	// wait for ^C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	// setup cli option parser
	parser := argparse.NewParser("mpdcord", "Discord Rich Presence for MPD written in Go")

	// get config path
	configHome := os.Getenv("XDG_CONFIG_HOME")
	var defaultConfigPath string
	if configHome != "" {
		defaultConfigPath = configHome + "mpdcord.toml"
	} else {
		homePath, _ := os.UserHomeDir()
		defaultConfigPath = homePath + "/.config/mpdcord.toml"
	}

	configPath := parser.String("c", "config", &argparse.Options{
		Required: false,
		Default:  defaultConfigPath,
        Help: "Specify non-standard config path.",
	})
	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Default:  false,
        Help: "Output additional information, useful for debugging.",
	})

	argerr := parser.Parse(os.Args)
	if argerr != nil {
		fmt.Print(parser.Usage(argerr))
		panic(argerr)
	}

	// check config path
	configContent, err := os.ReadFile(*configPath)
	if err != nil {
		ui.Warn("Couldn't read config at " + *configPath + ", using default values")
	} else if *verbose {
		ui.Info("Contents of " + *configPath)
		fmt.Println(string(configContent))
	}
	// read TOML values from config
	var config Config
	if err := toml.Unmarshal(configContent, &config); err != nil {
		panic(err)
	}
	// merge with default config
	mergo.Merge(&config, DefaultConfig)
	// pretty print current config
	if *verbose {
		prettyConfig := new(bytes.Buffer)
		if err := toml.NewEncoder(prettyConfig).Encode(config); err != nil {
			panic(err)
		}
		ui.Info("Current config:")
		fmt.Println(prettyConfig.String())
	}

	// connect to MPD
	ui.Running("Connecting to MPD at " + config.Address + " using " + config.Network)
	conn, err := mpd.Dial(config.Network, config.Address)
	if err != nil {
		ui.Error("Failed to connect to MPD")
		panic(err)
	} else {
		ui.Success("Connected to MPD")
	}

	// login to discord
	ui.Running("Logging in to Discord as " + strconv.Itoa(config.ID))
	login := client.Login(strconv.Itoa(config.ID))
	if login != nil {
		ui.Error("Couldn't log in to Discord")
		panic(login)
	} else {
		ui.Success("Logged in to Discord")
	}

	// listen to MPD events
	watcher, err := mpd.NewWatcher(config.Network, config.Address, config.Password)

	go func() {
		for range watcher.Event {
			// get and possibly print current status
            song, err := conn.CurrentSong()
			if err != nil {
				ui.Error("Couldn't get current song")
				panic(err)
			}
            status, err := conn.Status()
			if err != nil {
				ui.Error("Couldn't get status")
				panic(err)
			}
            stats, err := conn.Stats()
			if err != nil {
				ui.Error("Couldn't get stats")
				panic(err)
			}

			if *verbose {
				outsong, _ := json.Marshal(song)
				outstatus, _ := json.Marshal(status)
				ui.Info("Current status:")
				fmt.Println(string(outstatus))
				ui.Info("Current song:")
				fmt.Println(string(outsong))
			}

			// format strings from config
			details, err := interpol.WithMap(config.Format.Details, FormatMap(status, song, stats))
            if err != nil {
                ui.Error("Invalid formatting:")
                fmt.Println(config.Format.Details)
                panic(err)
            } else if *verbose {
                ui.Info("Details:")
                fmt.Println(details)
            }

			state, err := interpol.WithMap(config.Format.State, FormatMap(status, song, stats))
            if err != nil {
                ui.Error("Invalid formatting:")
                fmt.Println(config.Format.State)
                panic(err)
            } else if *verbose {
                ui.Info("State:")
                fmt.Println(state)
            }

			// get time when current song finishes
			elapsed, _ := time.ParseDuration(status["elapsed"] + "s")
			duration, _ := time.ParseDuration(status["duration"] + "s")
			start := time.Now()

			// define activity for RPC
			var activity = client.Activity{
				Details: details,
				State:   state,
				Timestamps: &client.Timestamps{
					Start: &start,
				},
			}


			if config.Format.Remaining {
				end := start.Add(duration).Add(-elapsed)
				activity.Timestamps.End = &end
			}

			// no need to clutter up the terminal every time you pause
			if *verbose {
				out, _ := json.Marshal(activity)
				ui.Running("Setting RPC status")
				fmt.Println(string(out))
			}

			client.SetActivity(activity)
		}
	}()

	<-done
	ui.Running("Closing MPD connection")
	conn.Close()
	ui.Running("Logging out")
	client.Logout()
}
