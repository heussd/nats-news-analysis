package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/heussd/nats-news-analysis/pkg/fulltextrss"
	"github.com/stretchr/testify/assert"
)

func TestNoMoreHTML(t *testing.T) {
	var fulltext = fulltextrss.RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")
	fulltext.Content = fulltext.Content + " <p>"

	assert.Equal(t, true, strings.Contains(fulltext.Content, "<p>"))

	var merged = prepareAndCleanString(fulltext)

	fmt.Printf("Here's the clean string %s\n", merged)
	assert.Equal(t, false, strings.Contains(merged, "<p>"))
}
