package ngrams

import (
	"strings"
	"unicode"
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
		if j, ok := index[gram]; ok {
			ngram[j].Count++
			continue
		}
		index[gram] = len(ngram)
		ngram = append(ngram, NGram{Words: gram, Count: 1, NGram: n})
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
