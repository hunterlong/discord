package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"layeh.com/gopus"
	"os/exec"
	"strconv"
	"sync"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

var (
	speakers         map[uint32]*gopus.Decoder
	opusEncoder      *gopus.Encoder
	mu               sync.Mutex
	playableChannels []*YoutubeOut
	onIndex          int
	onChannel        int
)

func PlayAudioFile(v *discordgo.VoiceConnection, filename string, stop <-chan bool) {
	youtubeDl := exec.Command("youtube-dl", "-f", "251", filename, "-o", "-")
	youtubeOut, err := youtubeDl.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = youtubeDl.Start()
	if err != nil {
		panic(err)
	}

	ffmpegRun := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegRun.Stdin = youtubeOut
	ffmpegout, err := ffmpegRun.StdoutPipe()
	if err != nil {
		panic(err)
	}

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	// Starts the ffmpeg command
	err = ffmpegRun.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		<-stop
		err = ffmpegRun.Process.Kill()
		err = youtubeDl.Process.Kill()
	}()

	// Send "speaking" packet over the voice websocket
	err = v.Speaking(true)
	if err != nil {
		panic(err)
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := v.Speaking(false)
		if err != nil {
			panic(err)
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	closer := make(chan bool)
	go func() {
		SendPCM(v, send)
		closer <- true
	}()

	for {
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			panic(err)
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-closer:
			return
		}
	}
}

func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}
	var err error

	opusEncoder, err = gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		panic(err)
	}

	for {

		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			fmt.Println("song ended")
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			panic(err)
		}

		if v.Ready == false || v.OpusSend == nil {
			// OnError(fmt.Sprintf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			// Sending errors here might not be suited
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}
