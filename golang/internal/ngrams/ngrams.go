package ngrams

import (
	"strings"
	"unicode"

	"github.com/tsawler/prose/v3"
)

func generateNGrams(text string, n int) (ngram []NGram, err error) {
	if n < 1 {
		return ngram, err
	}

	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	index := make(map[string]int)
	for i := 0; i+n <= len(words); i++ {
		gram := strings.Join(words[i:i+n], " ")
		index[gram] = len(ngram)
		ngram = append(ngram, NGram{Words: gram, NGram: n})
	}

	return ngram, err
}

func GenerateNGramStatistics(text string, minimumNGramSize int, maximumNGramSize int) (ngrams []NGram, err error) {
	var ngram []NGram

	for n := minimumNGramSize; n <= maximumNGramSize; n++ {
		ngram, err = generateNGrams(text, n)
		if err != nil {
			return ngrams, err
		}
		ngrams = append(ngrams, ngram...)
	}

	return ngrams, nil
}

func ParseAndGenerateStatistics(text string, minimumNGramSize int, maximumNGramSize int) (ngrams []NGram, err error) {
	doc, err := prose.NewDocument(text, prose.WithExtraction(false))
	if err != nil {
		return nil, err
	}

	words := []string{}
	prevI := 0
	for i, tok := range doc.Tokens() {
		if strings.HasPrefix(tok.Tag, "N") || strings.HasPrefix(tok.Tag, "C") || strings.HasPrefix(tok.Tag, "J") {
			words = append(words, strings.ToLower(tok.Text))
			prevI = i
		}
		if prevI+1 == i {
			// At this point we have identified the longest continuous sequence of nouns, adjectives, and conjunctions.
			if len(words) > 0 {
				newNGrams, err := GenerateNGramStatistics(strings.Join(words, " "), minimumNGramSize, maximumNGramSize)
				if err != nil {
					return nil, err
				}
				ngrams = append(ngrams, newNGrams...)
				words = []string{}
			}
		}
	}

	// Handle any remaining words after the loop
	if len(words) > 0 {
		newNGrams, err := GenerateNGramStatistics(strings.Join(words, " "), minimumNGramSize, maximumNGramSize)
		if err != nil {
			return nil, err
		}
		ngrams = append(ngrams, newNGrams...)
	}

	return ngrams, nil
}
