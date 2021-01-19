package main

import (
	"log"
	"net/http"

	"github.com/laytan/ytcrawler/pkg/crawler"
)

func main() {
	http.HandleFunc("/", crawler.CrawlChannel)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
