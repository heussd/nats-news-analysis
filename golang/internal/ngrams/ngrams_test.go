package ngrams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate1Grams(t *testing.T) {
	text := "This is a test test string for generating 1-grams."
	ngrams, err := generateNGrams(text, 1)

	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	assert.Equal(t, 9, len(ngrams))
	first := ngrams[0]
	assert.Equal(t, "This", first.Words)
	assert.Equal(t, 1, first.Count)
	fourth := ngrams[3]
	assert.Equal(t, "test", fourth.Words)
	assert.Equal(t, 2, fourth.Count)
	fiveth := ngrams[4]
	assert.Equal(t, "string", fiveth.Words)
	assert.Equal(t, 1, fiveth.Count)
}

func TestGenerate2Grams(t *testing.T) {
	text := "This is a test test string for generating 2-grams."
	ngrams, err := generateNGrams(text, 2)

	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	assert.Equal(t, 9, len(ngrams))
	first := ngrams[0]
	assert.Equal(t, "This is", first.Words)
	assert.Equal(t, 1, first.Count)
	fourth := ngrams[3]
	assert.Equal(t, "test test", fourth.Words)
	assert.Equal(t, 1, fourth.Count)
}

func TestGenerate3Grams(t *testing.T) {
	text := "This is a test test string for generating 3-grams."
	ngrams, err := generateNGrams(text, 3)

	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	assert.Equal(t, 8, len(ngrams))
	first := ngrams[0]
	assert.Equal(t, "This is a", first.Words)
	assert.Equal(t, 1, first.Count)
	fourth := ngrams[3]
	assert.Equal(t, "test test string", fourth.Words)
	assert.Equal(t, 1, fourth.Count)
}

func TestGenerate4Grams(t *testing.T) {
	text := "This is a test test string for generating 4-grams."
	ngrams, err := generateNGrams(text, 4)

	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	assert.Equal(t, 7, len(ngrams))
	first := ngrams[0]
	assert.Equal(t, "This is a test", first.Words)
	assert.Equal(t, 1, first.Count)
	fourth := ngrams[4]
	assert.Equal(t, "test string for generating", fourth.Words)
	assert.Equal(t, 1, fourth.Count)
}
