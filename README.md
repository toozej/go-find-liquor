# go-find-liquor (GFL)

Oregon Liquor Search Notification Service using [the OLCC Liquor Search website](http://www.oregonliquorsearch.com/), Go, and the [nikoksr/notify library](https://github.com/nikoksr/notify).

## Features

- Search [Oregon Liquor Search](http://www.oregonliquorsearch.com/) for specific liquor items
- Search by product name or item code
- Configurable search radius based on zip code
- Automatic age verification handling
- Random user agent rotation to avoid detection
- Random delays between searches to simulate human behavior
- Multiple notification methods:
  - Gotify (direct API integration)
  - Slack
  - Telegram
  - Discord
  - Pushover
  - Pushbullet
- Configurable search interval
- One-time or continuous search mode

## Usage with Docker

Dockerized GFL uses a final image based on [Distroless Debian static](https://github.com/GoogleContainerTools/distroless) since it includes minimal dependencies but most crucially ca-certificates (which are needed to POST notifications to HTTPS endpoints).

```bash
cp config.example.yaml config.yaml
# edit config.yaml
make build run
```

### Docker Compose

```bash
cp config.example.yaml config.yaml
# edit config.yaml
make up
```

## Installation

```bash
# Clone the repository
git clone https://github.com/toozej/go-find-liquor.git
cd go-find-liquor

# Build the application
make local-build
```

## Configuration

GFL can be configured using a YAML configuration file, environment variables, or command-line flags.

### Configuration File

Create a `config.yaml` file in the same directory as the executable:

```bash
cp config.example.yaml out/config.yaml
```

### Environment Variables

All configuration options can be set using environment variables with the `GFL_` prefix:

```bash
export GFL_ITEMS="Blanton's,Eagle Rare,Buffalo Trace"
export GFL_ZIPCODE="97201"
export GFL_DISTANCE="15"
export GFL_INTERVAL="6h"
```

### Command-Line Flags

Basic options can be set using command-line flags:

```bash
# Run in debug mode
./out/go-find-liquor -d

# Use a specific config file
./out/go-find-liquor -c /path/to/config.yaml

# Run search once and exit
./out/go-find-liquor -o
```

## Usage Examples

### Run a single search and exit

```bash
make local-run

# or alternatively without using the provided Makefile
./out/go-find-liquor --once
```

### Run continuously with the default interval

```bash
./out/go-find-liquor
```

### Run with a specific config file

```bash
./out/go-find-liquor --config /path/to/config.yaml
```

### View version information

```bash
./out/go-find-liquor version
```

### Generate man pages

```bash
./out/go-find-liquor man --dir /usr/local/share/man/man1
```

## Notification Types

GFL supports multiple notification methods:

### Gotify

```yaml
notifications:
  - type: gotify
    endpoint: "https://gotify.example.com"
    credential:
      token: "YOUR_GOTIFY_TOKEN"
```

### Slack

```yaml
notifications:
  - type: slack
    credential:
      token: "YOUR_SLACK_TOKEN"
      channel_id: "https://exampleorg.slack.com/archives/XXXXXXXXXXXXXXXXXXXXXXXX"
```

### Telegram

```yaml
notifications:
  - type: telegram
    credential:
      token: "YOUR_TELEGRAM_BOT_TOKEN"
      chat_id: "YOUR_CHAT_ID"
```

### Discord

```yaml
notifications:
  - type: discord
    credential:
      webhook_url: "https://discord.com/api/webhooks/000000000000000000/XXXXXXXXXXXXXXXXXXXXX"
```

### Pushover

```yaml
notifications:
  - type: pushover
    credential:
      token: "YOUR_PUSHOVER_TOKEN"
      receipient_id: "XXXXXXXXXXXXX"
```

### Pushbullet

```yaml
notifications:
  - type: pushbullet
    credential:
      token: "YOUR_PUSHBULLET_TOKEN"
      device_nickname: "XXXXXXXXXXXXX"
```

## Background

go-find-liquor, or GFL for short, was built since it is increasingly difficult to find some liquors at Oregon liquor stores due to short supply, mis-management, antiquated technology, etc. GFL was born to make it easier to find just the right bottle. Also, fun fact, GFL's alternative name is "good-fucking-luck", as in good luck finding those rare bottles ;).

## changes required to update golang version
- run `./scripts/update_golang_version.sh $NEW_VERSION_GOES_HERE`
