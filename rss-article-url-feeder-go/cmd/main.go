package main

import (
	"fmt"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/feed"
	queue "github.com/heussd/nats-news-keyword-matcher.go/internal/nats"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/urls"
	"time"
)

func main() {

	for {
		for _, feedUrl := range urls.Urls {
			for _, articleUrl := range feed.ArticleUrls(feedUrl) {
				queue.Publish(articleUrl)
			}
		}

		fmt.Printf("Waiting %d hour to retry...\n", 1)
		time.Sleep(1 * time.Hour)
		urls.ReloadUrls()
	}
}
