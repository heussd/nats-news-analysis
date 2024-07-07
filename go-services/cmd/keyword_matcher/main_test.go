package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	queue "nats-news-analysis/internal/nats"
	"nats-news-analysis/pkg/fulltextrss"
)

func TestNoMoreHTML(t *testing.T) {
	var fulltext = fulltextrss.RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")
	fulltext.Content = fulltext.Content + " <p>"

	assert.Equal(t, true, strings.Contains(fulltext.Content, "<p>"))

	var merged = prepareAndCleanString(fulltext)

	fmt.Printf("Here's the clean string %s\n", merged)
	assert.Equal(t, false, strings.Contains(merged, "<p>"))
}

func TestAnalysisFromQueue(t *testing.T) {

	queue.PublishArticle("https://www.tagesschau.de/wissen/forschung/dinosaurier-anpassung-klima-100.html")

	url := ""
	queue.WithArticleUrls(func(m *nats.Msg) {
		url = string(m.Data)
	})

	assert.Equal(t, "https://www.tagesschau.de/wissen/forschung/dinosaurier-anpassung-klima-100.html", url)
}
