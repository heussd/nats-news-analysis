package feed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func fetchNoErr(t *testing.T, url string) []string {
	res, err := FetchFeedAndExtractArticleUrls(url)
	assert.NoError(t, err)
	return res
}

func TestLocalIT(t *testing.T) {
	assert.True(t, len(fetchNoErr(t, "https://www.tagesschau.de/xml/rss2/")) > 0)
	assert.True(t, len(fetchNoErr(t, "https://www.heise.de/rss/heise-atom.xml")) > 0)
	assert.True(t, len(fetchNoErr(t, "https://neverworkintheory.org/atom.xml")) > 0)
}

func TestError(t *testing.T) {
	_, err := FetchFeedAndExtractArticleUrls("hello world")
	assert.Error(t, err)
}

func TestBaseUrl(t *testing.T) {
	assert.Equal(t, "https://www.example.com", getBaseUrl("https://www.example.com/feed.xml"))
	assert.Equal(t, "https://user@pass:localhost:8080", getBaseUrl("https://user@pass:localhost:8080/user/1000/profile?p=n#abc"))
}
