package ngrams

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/heussd/nats-news-analysis/internal/htmlsanitise"
	"github.com/heussd/nats-news-analysis/internal/model"
	"github.com/tsawler/prose/v3"
)

var (
	// remove purely numeric phrases (digits and spaces only)
	reOnlyNumbers = regexp.MustCompile(`^[0-9]+(\s+[0-9]+)*$`)
	// remove common HTML entity leftovers (e.g. "39", "34", "gt", "lt")
	reHTMLEntities = regexp.MustCompile(`(^|\s)(39|34|gt|lt)(\s|$)`)
	// remove parser/wiki template artifacts
	reWikiArtifacts = regexp.MustCompile(`(^|\s)(parser|navbox|hlist|reflist|liststyle|mw|cs1)(\s|$)`)
	// remove style/template boilerplate
	reStyleBoilerplate = regexp.MustCompile(`font size|font weight|background color|output div|references list|list style|style type|not skin`)
	// remove CSS/DOM noise
	reCSSNoise = regexp.MustCompile(`none none|padding 0|first child|last child|child before|child after|html skin|skin theme|theme clientpref|output [a-z0-9_]+|doi [0-9]+|id lock`)
)

func isNoise(words string) bool {
	w := strings.ToLower(words)
	return reOnlyNumbers.MatchString(words) ||
		reHTMLEntities.MatchString(w) ||
		reWikiArtifacts.MatchString(w) ||
		reStyleBoilerplate.MatchString(w) ||
		reCSSNoise.MatchString(w)
}

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

func ParseAndGenerateStatistics(news *model.News, minimumNGramSize int, maximumNGramSize int) (ngrams []NGram, err error) {
	var text string
	text = news.Content
	text = htmlsanitise.Sanitize(text)

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

	// Filter out noise: purely numeric phrases, HTML entity leftovers,
	// wiki/parser artifacts, style boilerplate, and CSS/DOM noise.
	filtered := ngrams[:0]
	for _, ng := range ngrams {
		if !isNoise(ng.Words) {
			filtered = append(filtered, ng)
		}
	}
	ngrams = filtered

	return ngrams, nil
}
