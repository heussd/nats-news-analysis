package main

import (
	"fmt"

	"github.com/heussd/nats-news-analysis/internal/feed"
	"github.com/heussd/nats-news-analysis/internal/model"
	queue "github.com/heussd/nats-news-analysis/internal/nats"
	"github.com/heussd/nats-news-analysis/pkg/utils"
)

func main() {
	var (
		input    = queue.AddStreamOrDie(utils.GetEnv("NATS_INPUT_STREAM", "feed-urls"))
		output   = queue.AddStreamOrDie(utils.GetEnv("NATS_OUTPUT_STREAM", "article-urls"))
		consumer = queue.AddConsumerOrDie(input, utils.GetEnv("NATS_CONSUMER", "default"))
	)

	var err = queue.Subscribe(input, consumer,
		func(f *model.Feed) {
			feedUrl := f.Url

			var articleUrls []string
			var err error
			if articleUrls, err = feed.FetchFeedAndExtractArticleUrls(feedUrl); err != nil {
				fmt.Printf("error fetching %s: %w", feedUrl, err)
				return
			}

			fmt.Printf("Found %d articles in %s", len(articleUrls), feedUrl)
			for _, articleUrl := range articleUrls {
				queue.Publish(output,
					model.Article{
						Url: articleUrl,
					},
					articleUrl,
					true,
				)
			}
		}, true,
	)

	if err != nil {
		panic(err)
	}
}
