package ngrams

import (
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStopwords(t *testing.T) {
	stopwordsExact, stopwordsAny := CachedStopwords()

	assert.True(t, slices.Contains(stopwordsExact, "the"))
	assert.True(t, slices.Contains(stopwordsAny, "jan"))
}

func TestStopwordsRegEx(t *testing.T) {
	_, stopwordsAny := CachedStopwords()

	r := regexp.MustCompile(`(^|\s)(` + strings.Join(stopwordsAny, "|") + `)(\s|$)`)
	assert.True(t, r.MatchString("aug 8"))
}

func TestIsNoise(t *testing.T) {
	assert.True(t, isNoise("123 456"))
	assert.False(t, isNoise("the quick brown fox"))
	assert.True(t, isNoise("https://example.com"))
	assert.False(t, isNoise("hello world"))
	assert.True(t, isNoise("aug 5"))
	assert.True(t, isNoise("brien aug 8 most"))
}
