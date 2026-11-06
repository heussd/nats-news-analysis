package htmlsanitise

import (
	"strings"
	"testing"

	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSanitisation(t *testing.T) {
	news := &model.News{
		Title:   "Title",
		Content: "<p>HELLO WORLD</p>",
	}
	outcome := ExtractFields(news)
	outcome = Sanitize(outcome)

	assert.False(t, strings.Contains(outcome, "<p>"))
	assert.False(t, strings.Contains(outcome, "</p>"))
	assert.True(t, strings.Contains(outcome, "HELLO WORLD"))
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

func TestSanitiseCodeBlock(t *testing.T) {
	outcome := Sanitize("HELLO <pre><code>BEAUTIFUL</code></pre> WORLD")

	assert.False(t, strings.Contains(outcome, "<pre>"))
	assert.False(t, strings.Contains(outcome, "</pre>"))
	assert.False(t, strings.Contains(outcome, "<code>"))
	assert.False(t, strings.Contains(outcome, "</code>"))
	assert.Equal(t, "HELLO WORLD", outcome)
}
