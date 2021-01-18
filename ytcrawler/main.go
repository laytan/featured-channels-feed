package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

var (
	errChannelNotGiven = fmt.Errorf("Channel not given in url as ?channel={channel}")
	errCantUnmarshal   = func(err error) error {
		return fmt.Errorf("Error parsing Youtube data: %v", err)
	}
	errNoPromotedChannels = fmt.Errorf("This channel does not have any featured channels")
)

type result struct {
	Error    string     `json:"error"`
	Channels []*channel `json:"channels"`
}

type video struct {
	URL          string `json:"url"`
	PublishedAt  string `json:"publishedAt"`
	ThumbnailURL string `json:"thumbnail"`
	Title        string `json:"title"`
	Views        string `json:"views"`
}

type channel struct {
	SubscriberCount string  `json:"subscriberCount"`
	URLTitle        string  `json:"urlTitle"`
	DisplayTitle    string  `json:"displayTitle"`
	LatestVideos    []video `json:"latestVideos"`
}

func main() {
	http.HandleFunc("/", Function)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Function starts the program to return the latest videos from the featured channels of the given channel
func Function(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", 60*30))
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))

	channelName, err := extractChannelName(r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result{Error: err.Error()})
		return
	}

	channels, statusCode, err := getFeaturedChannels(channelName)
	if err != nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result{Error: err.Error()})
		return
	}

	if len(channels) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(result{Error: errNoPromotedChannels.Error()})
		return
	}

	json.NewEncoder(w).Encode(result{Channels: channels})
}

// extractChannelName parses the url for the first parameter as the name
func extractChannelName(URL *url.URL) (string, error) {
	for key, value := range URL.Query() {
		if key == "channel" && value[0] != "" {
			return value[0], nil
		}
	}

	return "", errChannelNotGiven
}

// getFeaturedChannels scrapes the youtube channel for its featured channels
// Note: Tried doing it with the yt API
// but there is currently a bug which makes the featured channels not be returned through the api
func getFeaturedChannels(channelName string) ([]*channel, int, error) {
	start := time.Now()
	channels := make([]*channel, 0)
	var err error
	statusCode := http.StatusOK

	geziyor.NewGeziyor(&geziyor.Options{
		Timeout: time.Second * 5,
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(fmt.Sprintf("https://youtube.com/user/%s/channels", channelName), g.Opt.ParseFunc)
		},
		ParseFunc: onChannelsPage(&channels),
		ErrorFunc: func(_ *geziyor.Geziyor, _ *client.Request, errr error) {
			fmt.Println(errr)
			err = errr
			statusCode = http.StatusInternalServerError
		},
	}).Start()

	fmt.Printf("GetFeaturedChannels took: %s for %d featured channels of %s\n", time.Since(start), len(channels), channelName)
	return channels, statusCode, err
}

// Ran on youtube.com/user/example/channels and runs onChannel for each featured channel on the page
func onChannelsPage(channels *[]*channel) func(*geziyor.Geziyor, *client.Response) {
	return func(g *geziyor.Geziyor, r *client.Response) {
		if r.HTMLDoc == nil {
			return
		}

		r.HTMLDoc.Find("a.ytd-grid-channel-renderer").Each(func(_ int, s *goquery.Selection) {
			url, exists := s.Attr("href")
			if !exists {
				log.Printf("%s does not contain a href attribute\n", s.Text())
				return
			}

			g.GetRendered(fmt.Sprintf("https://www.youtube.com%s/videos", url), onChannel(channels))
		})
	}
}

// Ran on youtube.com/c/example/videos and extracts the channel info and latest videos
func onChannel(channels *[]*channel) func(*geziyor.Geziyor, *client.Response) {
	return func(g *geziyor.Geziyor, r *client.Response) {
		if r.HTMLDoc == nil {
			return
		}

		channelTitle := r.HTMLDoc.Find("#inner-header-container .ytd-channel-name#text").Text()
		subscriberText := r.HTMLDoc.Find("#inner-header-container #subscriber-count").Text()
		urlTitle := strings.Replace(r.Request.URL.Path, "/user/", "", 1)
		urlTitle = strings.Replace(urlTitle, "/channel/", "", 1)
		urlTitle = strings.Replace(urlTitle, "/videos", "", 1)

		videos := make([]video, 0)
		r.HTMLDoc.Find("#primary #contents #items .ytd-grid-renderer").Each(func(i int, s *goquery.Selection) {
			if i > 8 {
				return
			}

			videoURL, _ := s.Find("#thumbnail").Attr("href")
			thumbnailURL, _ := s.Find("#thumbnail img").Attr("src")
			videos = append(videos, video{
				URL:          videoURL,
				ThumbnailURL: thumbnailURL,
				PublishedAt:  s.Find("#details #metadata-line .ytd-grid-video-renderer:nth-child(2)").Text(),
				Title:        s.Find("#details #video-title").Text(),
				Views:        s.Find("#details #metadata-line .ytd-grid-video-renderer:nth-child(1)").Text(),
			})
		})

		*channels = append(*channels, &channel{
			SubscriberCount: subscriberText,
			URLTitle:        urlTitle,
			DisplayTitle:    channelTitle,
			LatestVideos:    videos,
		})
	}
}
