package main

import (
	"github.com/heussd/rss-article-url-feeder.go/internal/feed"
	queue "github.com/heussd/rss-article-url-feeder.go/internal/nats"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	queue.WithFeedUrls(func(m *nats.Msg) {
		feedUrl := string(m.Data)

		for _, articleUrl := range feed.ArticleUrls(feedUrl) {
			queue.Publish(articleUrl)
		}
	})

	// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}
