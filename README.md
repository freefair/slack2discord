# slack2discord

This is a Slack and Discord bot which syncs a Slack and Discord channel, so that every message is visible in both. 

## Usage 
The easiest way would be to use the prebuild docker-container. Just start it with this command `docker run frisch12/slack2discord my-container -- [ARGUMENTS]`

### Command line arguments
**redis-host**: Hostname and port of the redis to use (used to cache the last transmitted messages and usernames in case of a container restart). Default: 127.0.0.1:6379

**redis-pw**: Password for the redis to connect to. Default: <empty>`

**redis-db**: Database to use in redis. Default: `10`

**slack-token**: Token for your Slack bot (you have to create the Slack app yourself. Needs at least permission to read and write to channels that the bot is assigned to). Default: `<empty>`

**slack-channel**: ID of the Slack channel to sync. Default: `<empty>`

**discord-token**: Token for your Discord bot (you have to create the Slack app yourself. Needs at least permission to read and write to channels that the bot is assigned to). Default: `<empty>`

**discord-channel**: ID of the Discord channel to sync. Default: `<empty>`

## Licence
This software is licensed unter MIT
