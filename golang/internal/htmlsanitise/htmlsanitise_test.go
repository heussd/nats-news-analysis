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
	outcome := ExtractFields(news)
	outcome = Sanitize(outcome)

	assert.False(t, strings.Contains(outcome, "<p>"))
	assert.False(t, strings.Contains(outcome, "</p>"))
	assert.True(t, strings.Contains(outcome, "HELLO WORLD"))
	assert.Equal(t, fmt.Sprintf("%s %s %s", news.Title, news.Excerpt, "HELLO WORLD"), outcome)
}

func TestSanitisationWithEmptyContentAndNoHTML(t *testing.T) {
	outcome := Sanitize("hellö <img src='SHOULDBEHIDDEN'>asdasda</img>world asdasdsad & &amp; javascript:alert('HI') link another link https://example.com")

	assert.Equal(t, "hellö asdasda world asdasdsad & & javascript:alert('HI') link another link https://example.com", outcome)
}

func TestSanitisationP(t *testing.T) {
	outcome := Sanitize("<p>HELLO WORLD</p>")

	assert.False(t, strings.Contains(outcome, "<p>"))
	assert.False(t, strings.Contains(outcome, "</p>"))
	assert.True(t, strings.Contains(outcome, "HELLO WORLD"))
}
