package htmlsanitise

import (
	"fmt"
	"strings"
	"testing"

	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSanitisation(t *testing.T) {
	news := &model.News{
		Title:   "Title",
		Excerpt: "Excerpt",
		Content: "<p>HELLO WORLD</p>",
	}
	outcome := PrepareAndCleanString(news)

	assert.False(t, strings.Contains(outcome, "<p>"))
	assert.False(t, strings.Contains(outcome, "</p>"))
	assert.True(t, strings.Contains(outcome, "HELLO WORLD"))
	assert.Equal(t, fmt.Sprintf("%s %s %s", news.Title, news.Excerpt, "HELLO WORLD"), outcome)
}
