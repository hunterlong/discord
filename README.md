# Discord Audio Bot
A Discord audio bot that accept a list of youtube channel ID's to be streamed into you Discord chatroom/channels.
This Discord bot is best ran in a Docker container, I've provided examples below.

##### `docker pull hunterlong/discord`

## Bot Authentication
You need to authentication your bot to enter the channel by visiting a URL like below:
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
