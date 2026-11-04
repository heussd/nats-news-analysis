package defuddle

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalIT1(t *testing.T) {
	url := "https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html"
	article, err := RetrieveFullText(url)

	assert.NoError(t, err)

	assert.True(t,
		strings.Contains(article.Content,
			"Schnellere Abläufe könnten ebenfalls dazu beitragen"))

	assert.Equal(t,
		"Erneuerbare Energien: Doppelte Ernte mit Agri-Photovoltaik",
		article.Title)

	assert.Equal(t,
		"Ein Verbund aus Forschern, Landwirten und Klimaschützern fordert, die Agri-Photovoltaik auszubauen. Denn Solarzellen über Äckern und Obstwiesen könnten die Energiewende voranbringen - und Bauern eine \"zweite Ernte\" bescheren. Von Lena Puttfarcken.",
		article.Description)

	assert.Equal(t,
		"tagesschau.de",
		article.Domain)

	assert.Equal(t,
		"de",
		article.Language)

	assert.Equal(t,
		"2022-12-13T13:45:39.660Z",
		article.DatePublished)

	assert.Equal(t,
		"tagesschau.de",
		article.Author)

	assert.True(t,
		strings.Contains(article.ContentRawText, "Erneuerbare Energien Doppelte Ernte mit Agri-Photovoltaik Stand: 13.12.2022 • 14:45 Uhr Ein Verbund aus Forschern, Landwirten und Klimaschützern"),
	)

	assert.Equal(t,
		"tagesschau.de",
		article.Site)

	assert.Equal(t,
		802,
		article.WordCount)

	assert.Equal(t,
		url,
		article.URL)
}

func TestHasHTML(t *testing.T) {
	fulltext, _ := RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")
	fulltext.Content = fulltext.Content + " <p>"

	assert.Equal(t, true, strings.Contains(fulltext.Content, "<p>"))
}
