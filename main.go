package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strconv"
	"strings"
)

var (
	client     *discordgo.Session
	vconn      *discordgo.VoiceConnection
	youtubeKey string
	discordKey string
	chans      string
	channel    string
	guildId    string
	limitVids  int64
	bitRate    int64
)

func init() {
	channel = os.Getenv("CHANNEL_ID")
	guildId = os.Getenv("GUILD_ID")
	youtubeKey = os.Getenv("YOUTUBE")
	discordKey = os.Getenv("DISCORD")
	chans = os.Getenv("CHANNELS")

	lim := os.Getenv("LIMIT")
	limitVids, _ = strconv.ParseInt(lim, 10, 64)
	if limitVids == 0 {
		limitVids = 20
	}

	bitr := os.Getenv("BITRATE")
	bitRate, _ = strconv.ParseInt(bitr, 10, 64)
	if bitRate == 0 {
		bitRate = 64
	}
}

func updateList() {
	for _, v := range strings.Split(chans, ",") {
		out, err := Channel(v)
		fmt.Printf("Found %v videos for Channel: %s\n", len(out.Items), v)
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

	defer Close()

	playNext(onChannel, onIndex)
}

func Close() {
	vconn.Close()
	client.Close()
}

func playNext(chanIndex, vidIndex int) {
	chanObj := playableChannels[chanIndex]
	video := chanObj.Items[vidIndex]
	if video == nil {
		fmt.Printf("video %v not found\n", vidIndex)
		updateList()
		chanObj = playableChannels[0]
		video = chanObj.Items[vidIndex]
	}
	fmt.Printf("Now streaming: %s - %s - https://www.youtube.com/watch?v=%s\n", video.Snippet.ChannelTitle, video.Snippet.Title, video.ID.VideoID)
	if err := client.UpdateStatus(0, fmt.Sprintf("%s - %s", video.Snippet.ChannelTitle, video.Snippet.Title)); err != nil {
		fmt.Println(err)
	}

	PlayAudioFile(vconn, fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID.VideoID), make(chan bool))

	fmt.Printf("playing next channel %v item %v\n", chanIndex, vidIndex)

	vidIndex++

	if vidIndex >= len(chanObj.Items) {
		vidIndex = 0
		chanIndex++
	}

	if chanIndex >= len(playableChannels) {
		fmt.Println("updating list from channel: ", chanIndex)
		updateList()
		chanIndex = 0
		vidIndex = 0
	}

	playNext(chanIndex, vidIndex)
}
