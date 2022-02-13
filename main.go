package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
	"github.com/fhs/gompd/mpd"
	"github.com/hugolgst/rich-go/client"
	"github.com/pterm/pterm"
)

func main() {
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

	// get config
	err := GetConfig()
	if err != nil {
		pterm.Error.Printfln("Couldn't read config: %s", err)
	}

	// setup cli option parser
	parser := argparse.NewParser("mpdcord", "Discord Rich Presence for MPD written in Go")

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

	// connect to MPD
	connect := func() mpd.Client {
		spinner, _ := pterm.DefaultSpinner.Start("Connecting to MPD at " + DefaultConfig.Address + " using " + DefaultConfig.Network)
		conn, err := mpd.DialAuthenticated(DefaultConfig.Network, DefaultConfig.Address, DefaultConfig.Password)
		if err != nil {
			spinner.Fail("Failed to connect to MPD")
		} else {
			spinner.Success("Connected to MPD")
		}

		return *conn
	}

	// login to Discord
	login := func() error {
		spinner, _ := pterm.DefaultSpinner.Start("Logging in to Discord as " + strconv.Itoa(DefaultConfig.ID))
		err := client.Login(strconv.Itoa(DefaultConfig.ID))
		if err != nil {
			spinner.Fail()
		} else {
			spinner.Success()
		}
		return err
	}

	// listen to MPD events
	watcher, _ := mpd.NewWatcher(DefaultConfig.Network, DefaultConfig.Address, DefaultConfig.Password, "")
	defer watcher.Close()

	// try to connect to MPD and Discord
	conn := connect()

	// pinging and reconnecting
	reconnect := func() {
		err := conn.Ping()
		if err != nil {
			spinner, _ := pterm.DefaultSpinner.Start("Reconnecting to MPD")
			conn = connect()
			spinner.Success()
		}
		discord := login()
		if discord != nil {
			login()
		}
	}

	// timeout, _ := time.ParseDuration(config.Timeout)
	var song, status, stats mpd.Attrs
	var mpdmap map[string]string

	go func() {
		for range watcher.Event {
			// we have to reconnect every once in a while so we don't get timed out
			// there's probably a better way of fixing this but i'm too lazy to debug things properly
			reconnect()
			{
				// get current status
				song, _ = conn.CurrentSong()
				status, _ = conn.Status()
				stats, _ = conn.Stats()

				if *verbose {
					pterm.DefaultHeader.Printfln("--- %s", time.Now().Format(time.UnixDate))
					pterm.Info.Printfln("%s\n%s\n%s", song, status, stats)
				}

				// merge mpd status maps
				mpdmap = MergeMaps(song, status, stats)

				if *verbose {
					out, _ := json.Marshal(mpdmap)
					pterm.Info.Println("Current status:\n" + string(out))
				}

				// define activity for RPC
				var activity client.Activity

				if !(DefaultConfig.Format.PlayingOnly && mpdmap["state"] != "play") {
					activity = client.Activity{
						Details:    Formatted(DefaultConfig.Format.Details, mpdmap),
						State:      Formatted(DefaultConfig.Format.State, mpdmap),
						LargeImage: "mpd",
						LargeText:  Formatted(DefaultConfig.Format.LargeText, mpdmap),
						SmallImage: mpdmap["state"],
						SmallText:  Formatted(DefaultConfig.Format.SmallText, mpdmap),
						Timestamps: &client.Timestamps{},
					}

					// properly format time
					if mpdmap["state"] == "play" {
						elapsed, _ := time.ParseDuration(status["elapsed"] + "s")
						start := time.Now().Add(-elapsed)
						activity.Timestamps.Start = &start

						if *verbose {
							pterm.Info.Printfln("Elapsed: %s\nStart time: %s", elapsed.String(), start.Format(time.UnixDate))
						}

						if DefaultConfig.Format.Remaining {
							duration, _ := time.ParseDuration(status["duration"] + "s")
							end := time.Now().Add(duration).Add(-elapsed)
							activity.Timestamps.End = &end

							if *verbose {
								pterm.Info.Printfln("Duration: %s\nEnd time: %s", duration.String(), end.Format(time.UnixDate))
							}
						}

					}

					if *verbose {
						out, _ := json.Marshal(activity)
						spinner, _ := pterm.DefaultSpinner.Start("Setting RPC status")
						pterm.Info.Println(string(out))
						spinner.Success()
					}

					client.SetActivity(activity)
				} else {
					if *verbose {
						pterm.Info.Println("Logging out")
					}
					client.Logout()
				}
			}
		}
	}()

	<-done
	spinner, _ := pterm.DefaultSpinner.Start("Closing MPD connection")
	conn.Close()
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Logging out")
	client.Logout()
	spinner.Success()
}
