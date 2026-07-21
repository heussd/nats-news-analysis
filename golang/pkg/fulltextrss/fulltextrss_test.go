package fulltextrss

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalIT1(t *testing.T) {
	article := RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")

	assert.Equal(t,
		"Erneuerbare Energien: Doppelte Ernte mit Agri-Photovoltaik",
		article.Title)

	assert.True(t,
		strings.HasPrefix(article.Excerpt, "Stand: 13.12.2022"))

	assert.Equal(t,
		"de",
		article.Language)
}

func TestLocalIT2(t *testing.T) {
	article := RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")

	text := []string{
		article.Date,
		article.Language,
		article.Title,
		article.Excerpt,
	}
	fmt.Println(strings.Join(text, " ~#~ "))
}

func TestHasHTML(t *testing.T) {
	fulltext := RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")
	fulltext.Content = fulltext.Content + " <p>"

	assert.Equal(t, true, strings.Contains(fulltext.Content, "<p>"))
}
