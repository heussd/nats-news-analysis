package feed

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalIT(t *testing.T) {
	ArticleUrls("https://www.tagesschau.de/xml/rss2/")
	ArticleUrls("https://www.heise.de/rss/heise-atom.xml")
	ArticleUrls("https://neverworkintheory.org/atom.xml")
}

func TestBaseUrl(t *testing.T) {
	assert.Equal(t, "https://www.example.com", getBaseUrl("https://www.example.com/feed.xml"))
	assert.Equal(t, "https://user@pass:localhost:8080", getBaseUrl("https://user@pass:localhost:8080/user/1000/profile?p=n#abc"))
}
