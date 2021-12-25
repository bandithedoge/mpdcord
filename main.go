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
		Help:     "Specify non-standard config path",
	})
	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Default:  false,
		Help:     "Output additional information, useful for debugging",
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
			// get current status
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

			// merge mpd status maps
			mpdmap := MergeMaps(song, status, stats)

			if *verbose {
				out, _ := json.Marshal(mpdmap)
				ui.Info("Current status:")
				fmt.Println(string(out))
			}

			// define activity for RPC
			var activity = client.Activity{
				Details:    Formatted(config.Format.Details, mpdmap),
				State:      Formatted(config.Format.State, mpdmap),
                LargeImage: "mpd",
                LargeText: Formatted(config.Format.LargeText, mpdmap),
				Timestamps: &client.Timestamps{},
			}

			// properly format time
			if mpdmap["state"] == "play" {
				elapsed, _ := time.ParseDuration(status["elapsed"] + "s")
				start := time.Now().Add(-elapsed)
				activity.Timestamps.Start = &start

				if config.Format.Remaining {
					duration, _ := time.ParseDuration(status["duration"] + "s")
					end := time.Now().Add(duration).Add(-elapsed)
					activity.Timestamps.End = &end
				}
			}

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
