package ngrams

import (
	"testing"

	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerate1Grams(t *testing.T) {
	text := "This is a test test string for generating 1-grams."
	ngrams, err := generateNGrams(text, 1)

	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	assert.Equal(t, 10, len(ngrams))
	first := ngrams[0]
	assert.Equal(t, "This", first.Words)

	fourth := ngrams[3]
	assert.Equal(t, "test", fourth.Words)

	fiveth := ngrams[4]
	assert.Equal(t, "test", fiveth.Words)
}

func TestGenerate2Grams(t *testing.T) {
	text := "This is a test test string for generating 2-grams."
	ngrams, err := generateNGrams(text, 2)

	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	assert.Equal(t, 9, len(ngrams))
	first := ngrams[0]
	assert.Equal(t, "This is", first.Words)

	fourth := ngrams[3]
	assert.Equal(t, "test test", fourth.Words)
}

func TestGenerate3Grams(t *testing.T) {
	text := "This is a test test string for generating 3-grams."
	ngrams, err := generateNGrams(text, 3)

	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	assert.Equal(t, 8, len(ngrams))
	first := ngrams[0]
	assert.Equal(t, "This is a", first.Words)

	fourth := ngrams[3]
	assert.Equal(t, "test test string", fourth.Words)
}

func TestGenerate4Grams(t *testing.T) {
	text := "This is a test test string for generating 4-grams."
	ngrams, err := generateNGrams(text, 4)

	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	assert.Equal(t, 7, len(ngrams))
	first := ngrams[0]
	assert.Equal(t, "This is a test", first.Words)

	fourth := ngrams[4]
	assert.Equal(t, "test string for generating", fourth.Words)
}

func TestParseAndGenerateStatisticsFilterOut(t *testing.T) {
	var news model.News
	news.Content = "This"
	ngrams, err := ParseAndGenerateStatistics(&news, 1, 3)
	assert.NoError(t, err)
	assert.Empty(t, ngrams)
}

func TestParseAndGenerateStatisticsSample1(t *testing.T) {
	var news model.News
	news.Content = "This is macOS Tahoe Beta"
	ngrams, err := ParseAndGenerateStatistics(&news, 1, 4)
	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	expected := []NGram{
		{Words: "macos", NGram: 1},
		{Words: "tahoe", NGram: 1},
		{Words: "beta", NGram: 1},
		{Words: "macos tahoe", NGram: 2},
		{Words: "tahoe beta", NGram: 2},
		{Words: "macos tahoe beta", NGram: 3},
	}
	assert.Equal(t, expected, ngrams)
}

func TestParseAndGenerateStatistics(t *testing.T) {
	var news model.News
	news.Content = "Go is an open-source programming language created at Google. Windows 11 is an Operating System! RETRIEVAL-Augmented Generation. RAG is so 2023, just as GPT 3.5."
	ngrams, err := ParseAndGenerateStatistics(&news, 1, 4)
	assert.NoError(t, err)
	assert.NotEmpty(t, ngrams)

	expected := []NGram{
		{Words: "go", NGram: 1},
		{Words: "open", NGram: 1},
		{Words: "source", NGram: 1},
		{Words: "programming", NGram: 1},
		{Words: "language", NGram: 1},
		{Words: "open source", NGram: 2},
		{Words: "source programming", NGram: 2},
		{Words: "programming language", NGram: 2},
		{Words: "open source programming", NGram: 3},
		{Words: "source programming language", NGram: 3},
		{Words: "open source programming language", NGram: 4},
		{Words: "google", NGram: 1},
		{Words: "windows", NGram: 1},
		{Words: "windows 11", NGram: 2},
		{Words: "operating", NGram: 1},
		{Words: "system", NGram: 1},
		{Words: "operating system", NGram: 2},
		{Words: "retrieval", NGram: 1},
		{Words: "augmented", NGram: 1},
		{Words: "generation", NGram: 1},
		{Words: "retrieval augmented", NGram: 2},
		{Words: "augmented generation", NGram: 2},
		{Words: "retrieval augmented generation", NGram: 3},
		{Words: "rag", NGram: 1},
		{Words: "gpt", NGram: 1},
		{Words: "gpt 3", NGram: 2},
		{Words: "gpt 3 5", NGram: 3},
	}
	assert.Equal(t, expected, ngrams)
}
