package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

const (
	token = "NzE3Njg4NDA5ODgzMjc5NDQw.Xtd94Q.2afdQqI3yFi59WbJJpxpWenTCBI"
	channel = "717663692887556127"
	appId = "717688409883279440"
	guildId = "710101581030490182"
)

var (
	client *discordgo.Session
	vconn *discordgo.VoiceConnection
	youtubeKey string
)

func init() {
	youtubeKey = os.Getenv("YOUTUBE")
}

func main() {
	var err error

	out, err := Channel("UCYxRlFDqcWM4y7FfpiAN3KQ")
	if err != nil {
		panic(err)
	}

	for _, i := range out.Items {
		playableVids = append(playableVids, i)
		fmt.Println(i.ID.VideoID)
	}

	client, err = discordgo.New("Bot " + token)
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

	playNext(onIndex)

	defer Close()
}

func Close() {
	vconn.Close()
	client.Close()
}

func playNext(i int) {
	video := playableVids[i]
	PlayAudioFile(vconn, fmt.Sprintf(fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.ID.VideoID)), make(chan bool))
	playNext(i+1)
}
