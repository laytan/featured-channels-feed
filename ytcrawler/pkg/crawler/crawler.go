package crawler

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
	"github.com/geziyor/geziyor/middleware"
)

var (
	errChannelNotGiven = fmt.Errorf("Channel not givenn in url as ?channel={channel}")
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

type setCookiesMiddleware struct {
	cookies []*http.Cookie
}

func (m *setCookiesMiddleware) ProcessRequest(r *client.Request) {
	fmt.Printf("Setting %d cookies\n", len(m.cookies))
	for _, cookie := range m.cookies {
		r.AddCookie(cookie)
	}
}

// CrawlChannel starts the program to return the latest videos from the featured channels of the given channel
func CrawlChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Cache responses for 30 minutes
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

	cookieMiddleware := &setCookiesMiddleware{}

	geziyor.NewGeziyor(&geziyor.Options{
		Timeout: time.Second * 5,
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(
				fmt.Sprintf("https://www.youtube.com/%s/channels", channelName),
				g.Opt.ParseFunc,
			)
		},
		ParseFunc: onChannelsPage(channelName, &channels, cookieMiddleware),
		ErrorFunc: func(_ *geziyor.Geziyor, _ *client.Request, errr error) {
			fmt.Println(errr)
			err = errr
			statusCode = http.StatusInternalServerError
		},
		RequestMiddlewares: []middleware.RequestProcessor{cookieMiddleware},
		URLRevisitEnabled:  true,
	}).Start()

	fmt.Printf(
		"GetFeaturedChannels took: %s for %d featured channels of %s\n",
		time.Since(start),
		len(channels),
		channelName,
	)
	return channels, statusCode, err
}

// Ran on youtube.com/user/example/channels and runs onChannel for each featured channel on the page
func onChannelsPage(
	hostChannel string,
	channels *[]*channel,
	cookieMiddleware *setCookiesMiddleware,
) func(*geziyor.Geziyor, *client.Response) {
	consentSeen := false
	return func(g *geziyor.Geziyor, r *client.Response) {
		if r.HTMLDoc == nil {
			return
		}

		if r.Request.URL.Host == "consent.youtube.com" {
			fmt.Println("Accepting consent")
			if consentSeen {
				fmt.Println("[ERROR] Have seen consent already, but still redirected to consent.")
				return
			}

			request := url.Values{}

			r.HTMLDoc.Find(".saveButtonContainer").
				Children().
				Last().
				Find("input[type=\"hidden\"]").
				Each(func(i int, s *goquery.Selection) {
					if name, ok := s.Attr("name"); ok {
						if value, ok := s.Attr("value"); ok {
							request.Set(name, value)
						}
					}
				})

			res, err := http.PostForm("https://consent.youtube.com/save", request)
			if err != nil {
				fmt.Errorf("[ERROR] accepting consent/cookies: %w\n", err)
			}

			fmt.Printf("Status code %s for consent submission\n", res.Status)

			cookieMiddleware.cookies = res.Cookies()

			// Refire the request.
			g.GetRendered(
				fmt.Sprintf("https://www.youtube.com/%s/channels", hostChannel),
				g.Opt.ParseFunc,
			)

			consentSeen = true
			return
		}

		// Get videos from entered channel
		g.GetRendered(
			fmt.Sprintf("https://www.youtube.com/%s/videos", hostChannel),
			onChannel(channels),
		)

		// Get videos from featured channels
		r.HTMLDoc.Find("a.ytd-grid-channel-renderer").Each(func(_ int, s *goquery.Selection) {
			url, exists := s.Attr("href")
			if !exists {
				log.Printf("%s does not contain a href attribute\n", s.Text())
				return
			}

			g.GetRendered(fmt.Sprintf("https://www.youtube.com%s", url), onChannel(channels))
		})
	}
}

// Ran on youtube.com/c/example/videos and extracts the channel info and latest videos
func onChannel(channels *[]*channel) func(*geziyor.Geziyor, *client.Response) {
	return func(g *geziyor.Geziyor, r *client.Response) {
		if r.HTMLDoc == nil {
			return
		}

		fmt.Printf("Parsing %s\n", r.Request.URL.Path)

		channelTitle := r.HTMLDoc.Find("#inner-header-container .ytd-channel-name#text").Text()
		subscriberText := r.HTMLDoc.Find("#inner-header-container #subscriber-count").Text()

		urlTitle := strings.Replace(r.Request.URL.Path, "/user/", "", 1)
		urlTitle = strings.Replace(urlTitle, "/channel/", "", 1)
		urlTitle = strings.Replace(urlTitle, "/videos", "", 1)

		videos := make([]video, 0)
		r.HTMLDoc.Find(".ytd-browse #contents #items").First().Children().
			Each(func(i int, s *goquery.Selection) {
				if i >= 8 {
					return
				}

				fmt.Printf("Found video: %s\n", s.Find("#details #video-title").Text())

				videoURL, _ := s.Find("#thumbnail").Attr("href")
				thumbnailURL, _ := s.Find("#thumbnail img").Attr("src")
				videos = append(videos, video{
					URL:          videoURL,
					ThumbnailURL: thumbnailURL,
					PublishedAt: s.Find("#details #metadata-line .ytd-grid-video-renderer:nth-child(2)").
						Text(),
					Title: s.Find("#details #video-title").Text(),
					Views: s.Find("#details #metadata-line .ytd-grid-video-renderer:nth-child(1)").
						Text(),
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
