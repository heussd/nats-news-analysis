package main

import (
	"fmt"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/keywords"
	queue "github.com/heussd/nats-news-keyword-matcher.go/internal/nats"
	"github.com/heussd/nats-news-keyword-matcher.go/pkg/fulltextrss"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {

	queue.WithArticleUrls(func(m *nats.Msg) {
		var url = string(m.Data)
		var startTime = time.Now()
		var fulltext = fulltextrss.RetrieveFullText(url)

		var text = strings.Join([]string{
			fulltext.Title,
			fulltext.Excerpt,
			fulltext.Content,
		}, " ")

		var match, matchingText = keywords.Match(text)
		var elapsedTime = time.Since(startTime)
		if match {
			queue.PushToPocket(url, matchingText)
			fmt.Printf("✅ %s (analysis took %s)\n", url, elapsedTime)
		} else {
			fmt.Printf("❌ %s (analysis took %s)\n", url, elapsedTime)
		}
	})
	fmt.Println("\n🚀Keyword Matcher is ready to perform 🚀")

	// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}
