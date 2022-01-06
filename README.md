# mpdcord

![bangin](assets/screenshot.png)

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

## Installation

<!-- You can useBuilds mpdcord locally and places the binary at . You need to [install Go]() for this to work. -->

You can install `mpdcord` on any distro/OS by building it yourself with [Go](https://go.dev/doc/install). This will install the binary to `$GOPATH/bin/mpdcord` (defaults to `~/go/bin/mpdcord`), so make sure it's in your `$PATH`.

```bash
# latest stable
go install github.com/bandithedoge/mpdcord

# latest unstable
go install github.com/bandithedoge/mpdcord@latest
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

[Format]
  # First line
  Details = "{title}"
  # Second line
  State = "{artist}"
  # Text to display above the large image
  LargeText = "{album}"
  # Text to display above the small image
  SmallText = "{state}"
  # Time display type:
  #   - true: "XX:XX left"
  #   - false: "XX:XX elapsed"
  Remaining = false
```

## TODO

- [ ] Don't reconnect to Discord at every status change
- [ ] Figure out dynamically changing album covers (will definitely require a custom app ID)
