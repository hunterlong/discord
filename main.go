package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

const (
	appId = "717688409883279440"
)

var (
	client     *discordgo.Session
	vconn      *discordgo.VoiceConnection
	youtubeKey string
	discordKey string
	chans      string
	channel    string
	guildId    string
)

func init() {
	channel = os.Getenv("CHANNEL_ID")
	guildId = os.Getenv("GUILD_ID")
	youtubeKey = os.Getenv("YOUTUBE")
	discordKey = os.Getenv("DISCORD")
	chans = os.Getenv("CHANNELS")
}

func updateList() {
	for _, v := range strings.Split(chans, ",") {
		out, err := Channel(v)
		if err != nil {
			panic(err)
		}
		playableChannels = append(playableChannels, out)
	}
}

func main() {
	updateList()

	var err error
	client, err = discordgo.New("Bot " + discordKey)
	if err != nil {
		panic(err)
	}

	if err = client.Open(); err != nil {
		panic(err)
	}

	vconn, err = client.ChannelVoiceJoin(guildId, channel, false, true)
	if err != nil {
		panic(err)
	}

	playNext(onChannel, onIndex)

	defer Close()
}

func Close() {
	vconn.Close()
	client.Close()
}

func playNext(chanIndex, vidIndex int) {
	chanObj := playableChannels[chanIndex]
	video := chanObj.Items[vidIndex]
	fmt.Printf("Now streaming: %s - %s\n", video.Snippet.ChannelTitle, video.Snippet.Title)
	client.UpdateStatus(0, fmt.Sprintf("%s - %s", video.Snippet.ChannelTitle, video.Snippet.Title))
	PlayAudioFile(vconn, fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID.VideoID), make(chan bool))

	if vidIndex >= 3 {
		chanIndex++
		if chanIndex >= len(playableChannels) {
			chanIndex = 0
		} else {
			chanIndex++
			vidIndex = 0
		}
	} else {
		vidIndex += 1
	}

	if vidIndex >= 15 {
		updateList()
	}

	playNext(chanIndex, vidIndex)
}
