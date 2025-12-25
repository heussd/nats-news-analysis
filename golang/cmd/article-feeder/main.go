package main

import (
	"fmt"

	"github.com/heussd/nats-news-analysis/internal/feed"
	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
)

var publishSubject = "article-urls"

func main() {
	if err := queue.Subscribe(
		func(f *model.Feed) {
			feedURL := f.Url

			var articleUrls []string
			var err error
			if articleUrls, err = feed.FetchFeedAndExtractArticleUrls(feedURL); err != nil {
				fmt.Printf("error fetching %s: %v", feedURL, err)
				return
			}

			fmt.Printf("Found %d articles in %s\n", len(articleUrls), feedURL)
			for _, articleURL := range articleUrls {
				if _, err := queue.Publish(
					model.Article{
						Url: articleURL,
					},
					func(npo *queue.NatsPublishOpts) {
						npo.Subject = publishSubject
						npo.NatsMessageID = articleURL
						npo.PersistDeduplication = true
					},
				); err != nil {
					fmt.Printf("Failed to publish article URL %s: %v", articleURL, err)
				}
			}
		},
	); err != nil {
		panic(err)
	}
}
