package main

import (
	"fmt"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/keywords"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/model"
	queue "github.com/heussd/nats-news-keyword-matcher.go/internal/nats"
	"github.com/heussd/nats-news-keyword-matcher.go/pkg/fulltextrss"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var bm = bluemonday.StrictPolicy()

func prepareAndCleanString(fulltext fulltextrss.RSSFullTextResponse) string {
	var text = strings.Join([]string{
		fulltext.Title,
		fulltext.Excerpt,
		bm.Sanitize(fulltext.Content),
	}, " ")

	return text
}

func main() {

	queue.WithArticleUrls(func(m *nats.Msg) {
		var url = string(m.Data)
		var startTime = time.Now()
		var fulltext = fulltextrss.RetrieveFullText(url)
		var text = prepareAndCleanString(fulltext)

		var match, regexId = keywords.Match(text)
		var elapsedTime = time.Since(startTime)
		if match {
			queue.PushToPocket(model.Match{
				Url:     url,
				RegexId: regexId,
			})
			fmt.Printf("✅ %s (analysis took %s)\n", url, elapsedTime)
		} else {
			fmt.Printf("❌ %s (analysis took %s)\n", url, elapsedTime)
		}
	})
	fmt.Println("🚀Keyword Matcher is ready to perform 🚀")

	// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}
