package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Channel(channelId string) (*YoutubeOut, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?key=%s&channelId=%s&part=snippet,id&order=date&maxResults=%v", youtubeKey, channelId, limitVids)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		time.Sleep(15 * time.Second)
		return Channel(channelId)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code:", resp.StatusCode)
		time.Sleep(15 * time.Second)
		return Channel(channelId)
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Error:", err)
		time.Sleep(15 * time.Second)
		return Channel(channelId)
	}
	var vid *YoutubeOut
	err = json.Unmarshal(d, &vid)
	return vid, err
}

type YoutubeOut struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []*Item `json:"items"`
}

type Item struct {
	Kind    string   `json:"kind"`
	Etag    string   `json:"etag"`
	ID      *ItemID  `json:"id"`
	Snippet *Snippet `json:"snippet"`
}

type ItemID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type Snippet struct {
	PublishedAt time.Time `json:"publishedAt"`
	ChannelID   string    `json:"channelId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnails  struct {
		Default struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"default"`
		Medium struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"medium"`
		High struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"high"`
	} `json:"thumbnails"`
	ChannelTitle         string    `json:"channelTitle"`
	LiveBroadcastContent string    `json:"liveBroadcastContent"`
	PublishTime          time.Time `json:"publishTime"`
}
