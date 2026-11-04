package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSONOmitsContent(t *testing.T) {
	n := News{
		Content:        "<p>defuddle html</p>",
		Title:          "Headline",
		Description:    "Teaser",
		Domain:         "example.com",
		Language:       "en",
		DatePublished:  "2026-11-05",
		Author:         "Jane Doe",
		Site:           "Example",
		WordCount:      42,
		URL:            "https://example.com/article",
		ContentRawText: "defuddle html",
	}

	data, err := json.Marshal(n)
	assert.NoError(t, err)

	var decoded map[string]any
	assert.NoError(t, json.Unmarshal(data, &decoded))

	assert.NotContains(t, decoded, "content")
	assert.Equal(t, "Headline", decoded["title"])
	assert.Equal(t, "Teaser", decoded["description"])
	assert.Equal(t, "example.com", decoded["domain"])
	assert.Equal(t, "en", decoded["language"])
	assert.Equal(t, "2026-11-05", decoded["published"])
	assert.Equal(t, "Jane Doe", decoded["author"])
	assert.Equal(t, "Example", decoded["site"])
	assert.Equal(t, float64(42), decoded["wordCount"])
	assert.Equal(t, "https://example.com/article", decoded["url"])
	assert.Equal(t, "defuddle html", decoded["contentRawText"])
}

func TestMarshalJSONOmitsContentForPointerAndSlice(t *testing.T) {
	n := News{Content: "<p>defuddle html</p>", Title: "Headline"}

	data, err := json.Marshal(&n)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), `"content"`)

	data, err = json.Marshal([]News{n})
	assert.NoError(t, err)
	assert.NotContains(t, string(data), `"content"`)
}

func TestUnmarshalJSONKeepsContent(t *testing.T) {
	var n News
	err := json.Unmarshal([]byte(`{"content":"<p>defuddle html</p>","title":"Headline"}`), &n)
	assert.NoError(t, err)
	assert.Equal(t, "<p>defuddle html</p>", n.Content)
	assert.Equal(t, "Headline", n.Title)
}
