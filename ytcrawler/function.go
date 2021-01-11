package ytcrawler

// TODO: Rename package because we crawl and use the api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	errServer          = fmt.Errorf("Server error")
	errChannelNotGiven = fmt.Errorf("Channel not given in url as ?channel={channel}")
	errChannelNotFound = func(channel string) error {
		return fmt.Errorf("Channel: \"%s\" Not Found", channel)
	}
	errCollyErr = func(collyErr error) error {
		return fmt.Errorf("Error crawling Youtube: %v", collyErr)
	}
	errCantUnmarshal = func(err error) error {
		return fmt.Errorf("Error parsing Youtube data: %v", err)
	}
)

type result struct {
	Error  string                         `json:"error"`
	Result []*youtube.SearchResultSnippet `json:"result"`
}

// Function starts the program to return the latest videos from the featured channels of the given channel
func Function(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := godotenv.Load(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result{Error: errServer.Error()})
		return
	}

	channelName, err := extractChannelName(r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result{Error: err.Error()})
		return
	}

	// TODO: channels with multiple bars like EdSheeran don't work currently
	channels, statusCode, err := getFeaturedChannels(channelName)
	if err != nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result{Error: err.Error()})
		return
	}

	ytClient, err := getYTClient()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result{Error: err.Error()})
		return
	}

	videos := make([]*youtube.SearchResultSnippet, 0)
	resChan := make(chan *youtube.SearchResultSnippet)
	doneChan := make(chan bool)
	for _, chann := range channels {
		go func(chann channel) {
			call := ytClient.Search.List([]string{"snippet"})
			call.ChannelId(chann.ID)
			res, err := call.Do()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, res := range res.Items {
				resChan <- res.Snippet
			}
			doneChan <- true
		}(chann)
	}

	done := 0
Outer:
	for {
		select {
		case res := <-resChan:
			videos = append(videos, res)
		case <-doneChan:
			done++
			if done == len(channels) {
				break Outer
			}
		}
	}

	json.NewEncoder(w).Encode(result{Result: videos})
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
func getFeaturedChannels(channelName string) ([]channel, int, error) {
	channels := make([]channel, 0)
	var err error
	statusCode := http.StatusOK

	c := colly.NewCollector()

	// Check for error-page id which probably means the channel does not exist
	c.OnHTML("#error-page", func(e *colly.HTMLElement) {
		err = errChannelNotFound(channelName)
		statusCode = http.StatusNotFound
	})

	// Colly errors
	c.OnError(func(_ *colly.Response, e error) {
		err = errCollyErr(e)
		statusCode = http.StatusInternalServerError
	})

	// Get featured channels
	c.OnHTML("body", func(e *colly.HTMLElement) {
		rawYT, er := getRawYT(e)
		if er != nil {
			err = er
			statusCode = http.StatusInternalServerError
		}

		channels = extractChannels(rawYT)
	})

	c.Visit(fmt.Sprintf("https://youtube.com/c/%s/channels", channelName))
	c.Wait()
	return channels, statusCode, err
}

type ytdata struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents []struct {
								ItemSectionRenderer struct {
									Contents []struct {
										GridRenderer struct {
											Items []struct {
												GridChannelRenderer struct {
													ChannelID string
													Title     struct {
														SimpleText string
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func getRawYT(e *colly.HTMLElement) (*ytdata, error) {
	dat := e.ChildText("body > script:nth-child(11)")
	jsonData := dat[strings.Index(dat, "{") : strings.LastIndex(dat, "}")+1]
	data := &ytdata{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return data, errCantUnmarshal(err)
	}
	return data, nil
}

// Channel is a youtube channel
type channel struct {
	ID    string
	Title string
}

// extractChannels turns the complex ytdata object into a slice of channels
func extractChannels(data *ytdata) []channel {
	channels := make([]channel, 0)
	for _, a := range data.Contents.TwoColumnBrowseResultsRenderer.Tabs {
		for _, b := range a.TabRenderer.Content.SectionListRenderer.Contents {
			for _, c := range b.ItemSectionRenderer.Contents {
				for _, d := range c.GridRenderer.Items {
					channels = append(channels, channel{
						ID:    d.GridChannelRenderer.ChannelID,
						Title: d.GridChannelRenderer.Title.SimpleText,
					})
				}
			}
		}
	}
	return channels
}

// getYTClient sets up the youtube library
func getYTClient() (*youtube.Service, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		log.Println("No YOUTUBE_API_KEY in env")
		return &youtube.Service{}, errServer
	}

	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Error getting youtube service: %v\n", err)
		return &youtube.Service{}, errServer
	}

	return youtubeService, nil
}
