package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"layeh.com/gopus"
	"os"
	"os/exec"
	"strconv"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

var (
	opusEncoder      *gopus.Encoder
	playableChannels []*YoutubeOut
	onIndex          int
	onChannel        int
)

// PlayAudioFile
func PlayAudioFile(v *discordgo.VoiceConnection, filename string, stop <-chan bool) {
	youtubeDl := exec.Command("youtube-dl", "--no-color", "--audio-format", "best", "--audio-format", "opus", filename, "-o", "-")
	youtubeOut, err := youtubeDl.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err = youtubeDl.Start(); err != nil {
		panic(err)
	}

	ffmpegRun := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegRun.Stdin = youtubeOut
	ffmpegRun.Stderr = os.Stderr
	ffmpegout, err := ffmpegRun.StdoutPipe()
	if err != nil {
		panic(err)
	}

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	if err := ffmpegRun.Start(); err != nil {
		panic(err)
	}

	go func() {
		<-stop
		err = youtubeDl.Process.Kill()
		err = ffmpegRun.Process.Kill()
	}()

	if err := v.Speaking(true); err != nil {
		panic(err)
	}

	defer func() {
		if err := v.Speaking(false); err != nil {
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

	opusEncoder.SetBitrate(int(bitRate) * 1000)

	for {
		recv, ok := <-pcm
		if !ok {
			fmt.Println("song ended, or error")
			return
		}

		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			panic(err)
		}

		if v.Ready == false || v.OpusSend == nil {
			fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend)
			return
		}
		v.OpusSend <- opus
	}
}
