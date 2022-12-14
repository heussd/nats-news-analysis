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
)

func main() {

	queue.WithArticleUrls(func(m *nats.Msg) {
		var url = string(m.Data)
		fmt.Printf("Processing %s...\n", url)
		var fulltext = fulltextrss.RetrieveFullText(url)

		var text = strings.Join([]string{
			fulltext.Title,
			fulltext.Excerpt,
			fulltext.Content,
		}, " ")

		if keywords.Match(text) {
			queue.PushToPocket(url)
		}
	})
	fmt.Println("\n🚀Keyword Matcher is ready to perform 🚀")

	// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
}
