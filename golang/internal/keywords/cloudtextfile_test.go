package keywords

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var keywords []string
	var err error
	if keywords, err = RetrieveKeywordsFile(); err != nil {
		t.Error(err)
	}

	assert.Equal(t,
		12,
		len(keywords))

	assert.Equal(t,
		"(?i)\\b(Apple|peach)",
		keywords[0])
}

func TestCache(t *testing.T) {
	// First call to populate the cache
	keywords, err := CachedParsedKeywords()
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, keywords)
	assert.Greater(t, len(keywords), 0)

	// Store the time of the first cache generation
	firstGenerated := lastGenerated

	// Wait for a short duration and call again
	time.Sleep(1 * time.Second)
	keywords, err = CachedParsedKeywords()
	if err != nil {
		t.Error(err)
	}

	// Ensure the cache is still valid and hasn't been regenerated
	assert.Equal(t, firstGenerated, lastGenerated)
	assert.NotNil(t, keywords)
	assert.Greater(t, len(keywords), 0)

}

func TestHumanReadable(t *testing.T) {
	assert.Equal(t, "delicious pie recipes", humanReadable("(?i)(delicious).*(pie|recipes)"))

}
