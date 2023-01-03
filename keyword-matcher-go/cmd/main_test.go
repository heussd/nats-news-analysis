package main

import (
	"fmt"
	"github.com/heussd/nats-news-keyword-matcher.go/pkg/fulltextrss"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNoMoreHTML(t *testing.T) {
	var fulltext = fulltextrss.RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")

	assert.Equal(t, true, strings.Contains(fulltext.Content, "<p>"))

	var merged = prepareAndCleanString(fulltext)

	fmt.Printf("Here's the clean string %s\n", merged)
	assert.Equal(t, false, strings.Contains(merged, "<p>"))
}
