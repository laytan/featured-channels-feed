package ytcrawler

// TODO: Rename package because we crawl and use the api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

var (
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

// Function starts the program to return the latest videos from the featured channels of the given channel
func Function(w http.ResponseWriter, r *http.Request) {
	channelName, err := extractChannelName(r.URL)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	channels, err := getFeaturedChannels(channelName)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// TODO: ask yt api for latest videos of these channels
	// TODO: return the videos and errors in JSON

	fmt.Fprint(w, channels)
}

// extractChannelName parses the url for the first parameter as the name
func extractChannelName(URL *url.URL) (string, error) {
	for key, value := range URL.Query() {
		if key == "channel" {
			return value[0], nil
		}
	}

	return "", errChannelNotGiven
}

// getFeaturedChannels scrapes the youtube channel for its featured channels
// Note: Tried doing it with the yt API
// but there is currently a bug which makes the featured channels not be returned through the api
func getFeaturedChannels(channelName string) ([]channel, error) {
	channels := make([]channel, 0)
	var err error

	c := colly.NewCollector()

	// Check for error-page id which probably means the channel does not exist
	c.OnHTML("#error-page", func(e *colly.HTMLElement) {
		err = errChannelNotFound(channelName)
	})

	// Colly errors
	c.OnError(func(_ *colly.Response, e error) {
		err = errCollyErr(e)
	})

	// Get featured channels
	c.OnHTML("body", func(e *colly.HTMLElement) {
		rawYT, er := getRawYT(e)
		if er != nil {
			err = er
		}

		channels = extractChannels(rawYT)
	})

	c.Visit(fmt.Sprintf("https://youtube.com/c/%s/channels", channelName))
	c.Wait()
	return channels, err
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
