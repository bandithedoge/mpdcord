# mpdcord

## Installation

Builds mpdcord locally and places the binary at `$GOPATH/bin/mpdcord` (defaults to `~/go/bin/mpdcord`). You need to [https://go.dev/doc/install](install Go) for this to work.

```bash
# latest stable
go install github.com/bandithedoge/mpdcord

# latest unstable
go install github.com/bandithedoge/mpdcord@latest
```

## Usage

```console
❯ go run main
~ Connecting to MPD at localhost:6600 using tcp
✓ Connected to MPD
~ Logging in to Discord as 922175995828654100
✓ Logged in to Discord
^C
interrupt
~ Closing MPD connection
~ Logging out
```

## Configuration

Everything is configured in a TOML file located at `$XDG_CONFIG_HOME/mpdcord.toml` (defaults to `~/.config/mpdcord.toml`). Here is an example configuration populated with default values:

```toml
# Note: keys are not case-sensitive.

# Discord API application ID, use this to customize title and images
ID = 922175995828654100
# Where to connect to MPD
Address = "localhost:6600"
# How to connect to MPD
Network = "tcp"
# Optional MPD password
Password = ""

# All the formatting is done using values wrapped in curly braces, for example "{title}"
# Possible values: 
#   "volume",
#   "repeat",
#   "random",
#   "single",
#   "playlistlength",
#   "consume",
#   "audio",
#   "bitrate",
#   "album",
#   "artist",
#   "albumartist",
#   "composer",
#   "conductor",
#   "date",
#   "disc",
#   "ensemble",
#   "genre",
#   "grouping",
#   "label",
#   "location",
#   "movement",
#   "movementnumber",
#   "originaldate",
#   "performer",
#   "title",
#   "track",
#   "work",
#   "artists",
#   "albums",
#   "songs",

[Format]
  # First line
  Details = "{title}"
  # Second line
  State = "{artist}"
  # Time display type:
  #   - true: "XX:XX left"
  #   - false: "XX:XX elapsed"
  Remaining = false
```
