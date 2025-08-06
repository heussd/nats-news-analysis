package fulltextrss

import (
	"fmt"
	"strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestLocalIT1(t *testing.T) {
	var article = RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")

	assert.Equal(t,
		"Erneuerbare Energien: Doppelte Ernte mit Agri-Photovoltaik",
		article.Title)

	assert.Equal(t,
		"Stand: 13.12.2022 14:45 Uhr Ein Verbund aus Forschern, Landwirten und Klimaschützern fordert, die Agri-Photovoltaik auszubauen. Denn Solarzellen über Äckern und Obstwiesen könnten die Energiewende voranbringen - und Bauern eine \"zweite Ernte\" bescheren.",
		article.Excerpt)

	assert.Equal(t,
		"de",
		article.Language)
}

func TestLocalIT2(t *testing.T) {
	var article = RetrieveFullText("https://www.tagesschau.de/wissen/technologie/agri-photovoltaik-103.html")

	var text = []string{
		article.Date,
		article.Language,
		article.Title,
		article.Excerpt,
	}
	fmt.Println(strings.Join(text, " ~#~ "))
}
