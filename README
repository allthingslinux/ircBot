# Discord IRC Bot

A Discord bot that provides IRC management capabilities through Discord commands. Originally designed for [irc.atl.chat](https://atl.chat) and [All Things Linux](https://allthingslinux.org), but configurable for other IRC networks.

## Features

- **User Registration**: Register new IRC users through NickServ
- **User Deletion**: Delete IRC users with confirmation codes
- **Help System**: Display available commands
- **Graceful Shutdown**: Proper cleanup on termination signals
- **Error Handling**: Comprehensive error handling and logging
- **Timeout Protection**: Commands have built-in timeouts to prevent hanging

## Requirements

- [Go 1.24+](https://go.dev)
- Access to a Discord bot and its token
- An IRC account with operator privileges

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/allthingslinux/ircBot
   cd ircBot
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Copy the example environment file and configure it:
   ```bash
   cp .example.env .env
   ```

4. Fill out the `.env` file with your configuration (see Configuration section)

5. Build the application:
   ```bash
   go build -o ircBot
   ```

6. Run the bot:
   ```bash
   ./ircBot
   ```

## Available Commands

- `!help` - Display available commands

## Development

### Project Structure

```
├── main.go              # Application entry point and configuration
├── commands/
│   ├── commands.go      # Command routing and handler registration
│   └── utils.go         # Shared utilities
├── go.mod               # Go module definition
└── .env                 # Environment configuration
```

### Adding New Commands

1. Create a new file in the `commands/` directory
2. Implement your command handler with this signature:
   ```go
   func YourCommand(session *discordgo.Session, message *discordgo.MessageCreate, irccon *irc.Connection)
   ```
3. Register the command in the `init()` function:
   ```go
   func init() {
       CommandMap["yourcommand"] = YourCommand
   }
   ```

## Contributing

Pull requests and issues are welcome. This is a community project maintained by volunteers, so please be patient with response times.