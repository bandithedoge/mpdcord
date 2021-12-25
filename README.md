# mpdcord

## Installation

Builds mpdcord locally and places the binary at `$GOPATH/bin/mpdcord` (defaults to `~/go/bin/mpdcord`). You need to [install Go](https://go.dev/doc/install) for this to work.

```bash
# latest stable
go install github.com/bandithedoge/mpdcord

# latest unstable
go install github.com/bandithedoge/mpdcord@latest
```

## Usage

```console
usage: mpdcord [-h|--help] [-c|--config "<value>"] [-v|--verbose]

               Discord Rich Presence for MPD written in Go

Arguments:

  -h  --help     Print help information
  -c  --config   Specify non-standard config path. Default:
                 $XDG_CONFIG_HOME/mpdcord.toml
  -v  --verbose  Output additional information, useful for debugging. Default:
                 false
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
#   "album",
#   "albumartist",
#   "albums",
#   "artist",
#   "artists",
#   "audio",
#   "bitrate",
#   "composer",
#   "conductor",
#   "consume",
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
#   "playlistlength",
#   "random",
#   "repeat",
#   "single",
#   "songs",
#   "title",
#   "track",
#   "volume",
#   "work"

[Format]
  # First line
  Details = "{title}"
  # Second line
  State = "{artist}"
  # Text to display when hovering over the large image
  LargeText = "{album}"
  # Time display type:
  #   - true: "XX:XX left"
  #   - false: "XX:XX elapsed"
  Remaining = false
```
