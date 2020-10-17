[![Actions Status](https://github.com/maliur/sodaville/workflows/build/badge.svg)](https://github.com/maliur/sodaville/actions) [![Actions Status](https://github.com/maliur/sodaville/workflows/tests/badge.svg)](https://github.com/maliur/sodaville/actions)
```
            ███████╗ ██████╗ ██████╗  █████╗ ██╗   ██╗██╗██╗     ██╗     ███████╗
            ██╔════╝██╔═══██╗██╔══██╗██╔══██╗██║   ██║██║██║     ██║     ██╔════╝
            ███████╗██║   ██║██║  ██║███████║██║   ██║██║██║     ██║     █████╗  
            ╚════██║██║   ██║██║  ██║██╔══██║╚██╗ ██╔╝██║██║     ██║     ██╔══╝  
            ███████║╚██████╔╝██████╔╝██║  ██║ ╚████╔╝ ██║███████╗███████╗███████╗
            ╚══════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝  ╚═══╝  ╚═╝╚══════╝╚══════╝╚══════╝
```

Another Twitch chat bot

## Installation
### Requirements
- Go language install with your package manager or from [golang.org](https://golang.org/)
- OAuth token from twitch **Might show sensitive information** [Twitch chat OAuth generator](https://twitchapps.com/tmi/)
- SQLite3 install with your package manager or from [sqlite.org](https://sqlite.org/index.html)
- (Optional) Install Make

## Usage
Make sure that you have set these environment variables
```
OAUTH_TOKEN=<value>
BOT_USERNAME=<value-here>
CHANNEL_NAME=<value-here>
```

Run the database migration
```bash
sqlite3 sodaville.db < migrations.sql
```

To start the bot:
```bash
go build -o sodaville cmd/bot/main.go && ./sodaville
```
Or with Make:
```bash
make run
```


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
Released under [MIT license](https://raw.githubusercontent.com/maliur/sodaville/master/LICENSE)