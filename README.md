# Discord Audio Bot
A Discord audio bot that accept a list of youtube channel ID's to be streamed into you Discord chatroom/channels.
This Discord bot is best ran in a Docker container, I've provided examples below.

##### `docker pull hunterlong/discord`

## Installation
1. Create a new Application on discord by going to: [https://discord.com/developers/applications/](https://discord.com/developers/applications/).
2. Click on the Bot tab on the left, and Add a Bot.
3. Copy the Bot's Token Key, this will be set as `DISCORD` environment variable.
4. Confirm `SERVER MEMBERS INTENT` is switched ON for the bot.
5. Copy your Bot's Client ID and visit the URL: `https://discord.com/oauth2/authorize?scope=bot&client_id=<CLIENT_ID_HERE>`. Be sure to replace `<CLIENT_ID_HERE>` with your bot's client ID. This will authorize your bot to become a member of your channel.
6. Enable Developer Mode in the Discord app, and right click your Voice Chanel for the bot and click Copy ID, this will be used for `CHANNEL_ID`.
7. Right click your organization in Discord, and copy the Copy ID, this will be used for `GUILD_ID`.
8. Create a Youtube v3 API Key by going to [Google API Console](https://developers.google.com/youtube/registering_an_application), create a new app, the Youtube API key will used for `YOUTUBE`.
9. Collect some Youtube Channel ID numbers from your favorite places. For example, Channel ID [UCBOqkAGTtzZVmKvY4SwdZ2g](https://www.youtube.com/channel/UCBOqkAGTtzZVmKvY4SwdZ2g).

## Bot Authentication
You need to authenticate your bot to enter the channel by visiting a URL like below:
```
https://discord.com/oauth2/authorize?scope=bot&client_id=<CLIENT_ID>
```
Replace `<CLIENT_ID>` with the bot's client ID number and authorize the bot.

## Docker Environment Variables
- `YOUTUBE`: Found on Youtube's Credentials API page
- `DISCORD`: Discount bot client secret found on Discord's bot website
- `CHANNELS`: A comma delimited list of Youtube Channels
- `CHANNEL_ID`: The Discord chatroom/channel ID number
- `GUILD_ID`: The Discord organization ID number

## Docker Compose Example
```yaml
discord_music:
  container_name: discord_music
  image: hunterlong/discord
  restart: always
  environment:
    YOUTUBE: YOUTUBE_API_KEY
    DISCORD: DISCORD_BOT_SECRET
    CHANNELS: "UCw49uOTAJjGUdoAeUcp7tOg,UCw49uOTAJjGUdoAeUcp7tOg"
    CHANNEL_ID: 7101030EXAMPLE735568
    GUILD_ID: 7101EXAMPLE490182
  volumes:
    - discordaudio:/downloads:rw
```

## Packages Used
- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
- [bwmarrin/dgvoice](https://github.com/bwmarrin/dgvoice)
- [qmcgaw/youtube-dl-alpine](https://hub.docker.com/r/qmcgaw/youtube-dl-alpine)

## License
MIT. PR's accepted if you want to.
